package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Repo struct {
	SvnUrl string `json:"svn_url"`
	Json   *json.RawMessage
}

type PrvRepoUrls []struct {
	HtmlUrl string `json:"html_url"`
}

func main() {
	org := flag.String("org", "skycoin", "github organization name")
	out := flag.String("dir", ".", "destination dir where downloads will be stored")
	// filename := flag.String("f", "", "get urls from file. If not provided, https://api.github.com/orgs/:org will be reached to get all repos urls of an organization")
	token := flag.String("token", "", "private org access token")
	urlsOnly := flag.Bool("urls-only", false, "get and print repo urls only")

	flag.Parse()

	isPrivate, err := isPrivateOrg(*org)
	if err != nil {
		log.Fatalf("failed to check if org: %s is private, err: %v\n", *org, err)
		return
	}

	if isPrivate && *token == "" {
		log.Fatal("access token is empty, please set it with -token flag")
		return
	}

	// Check if output dir exists, create it if does not exist
	outDir := fmt.Sprintf("%s/%s", *out, *org)
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		// create the dir
		if err := os.MkdirAll(outDir, os.FileMode(0700)); err != nil {
			log.Fatalf("create output dir failed: %v", err)
			return
		}
	}

	var urls []string
	if isPrivate {
		urls, err = getPrivateRepoUrls(*org, *token)
		if err != nil {
			log.Fatal(err)
			return
		}
	} else {
		urls, err = getRepoUrls(*org)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	// }

	if *urlsOnly {
		for _, url := range urls {
			fmt.Println(url)
		}
		return
	}

	if err := downloadRepos(urls, *org, outDir); err != nil {
		log.Fatal(err)
	}
	return
}

func getTotalPages(link string) (string, int, error) {
	fmt.Println(link)
	links := strings.Split(link, ",")
	for _, l := range links {
		ll := strings.Split(l, ";")
		rel := strings.Trim(ll[1], " ")
		if rel == "rel=\"last\"" {
			url := strings.Trim(ll[0], " ")
			url = strings.TrimLeft(url, "<")
			url = strings.TrimRight(url, ">")
			surl := strings.Split(url, "=")
			baseURL := surl[0]
			pageStr := surl[1]
			page, err := strconv.Atoi(pageStr)
			if err != nil {
				return "", 0, err
			}

			return baseURL, page, nil
		}
	}
	return "", 0, errors.New("could not find last page")
}

func getPrivateRepoUrls(org, accessToken string) ([]string, error) {
	page := 1
	var urls []string
	for {
		url := fmt.Sprintf("https://api.github.com/orgs/%s/repos?access_token=%s&page=%d&per_page=100", org, accessToken, page)
		rsp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer rsp.Body.Close()

		prs := PrvRepoUrls{}
		if err := json.NewDecoder(rsp.Body).Decode(&prs); err != nil {
			return nil, err
		}

		if len(prs) == 0 {
			break
		}

		for _, url := range prs {
			urls = append(urls, url.HtmlUrl)
		}
		page++
	}
	return urls, nil
}

func getRepoUrls(org string) ([]string, error) {
	rsp, err := http.Get(fmt.Sprintf("https://api.github.com/orgs/%s/repos", org))
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		v, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		return nil, fmt.Errorf("fetch error: %v", string(v))
	}

	repos := []Repo{}
	if err := json.NewDecoder(rsp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	headLink := rsp.Header.Get("Link")
	if headLink == "" {
		var rs []string
		for _, r := range repos {
			rs = append(rs, r.SvnUrl)
		}
		return rs, nil
	}

	// fmt.Printf("Link: %+v\n", rsp.Header.Get("Link"))
	baseURL, totalPages, err := getTotalPages(rsp.Header.Get("Link"))
	if err != nil {
		return nil, err
	}

	if totalPages <= 1 {
		return nil, nil
	}

	for i := 2; i <= totalPages; i++ {
		func(page int) {
			r, err := http.Get(fmt.Sprintf("%s=%v", baseURL, page))
			if err != nil {
				log.Fatalf("failed to get page:%v, err: %v\n", i, err)
				return
			}
			defer r.Body.Close()

			var rs []Repo
			if err := json.NewDecoder(r.Body).Decode(&rs); err != nil {
				log.Fatal(err)
				return
			}
			repos = append(repos, rs...)
		}(i)
	}

	var urls []string
	for _, r := range repos {
		urls = append(urls, r.SvnUrl)
	}

	return urls, nil
}

func downloadRepos(urls []string, org, outDir string) error {
	fs, err := ioutil.ReadDir(outDir)
	if err != nil {
		return err
	}

	alreadyDown := map[string]bool{}

	for _, f := range fs {
		s := fmt.Sprintf("https://github.com/%s/%s", org, f.Name())
		alreadyDown[s] = true
	}

	for _, r := range urls {
		if alreadyDown[r] {
			continue
		}
		repoName := r[strings.LastIndex(r, "/")+1:]
		to := fmt.Sprintf("%s/%s", outDir, repoName)
		fmt.Println("start to clone:", r, "to", to)
		cmd := exec.Command("git", "clone", r, to)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}

		var errbuf bytes.Buffer
		cmd.Stderr = &errbuf
		cmd.Start()

		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
		}

		fmt.Println(errbuf.String())
		cmd.Wait()
		alreadyDown[r] = true
	}

	return nil
}

// check if the org is private
func isPrivateOrg(org string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos", org)
	rsp, err := http.Get(url)
	if err != nil {
		return false, err
	}

	defer rsp.Body.Close()
	repos := []Repo{}
	if err := json.NewDecoder(rsp.Body).Decode(&repos); err != nil {
		return false, err
	}
	return len(repos) == 0, nil
}

func getUrlsFromFile(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var urls []string
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		urls = append(urls, strings.Trim(scan.Text(), ""))
	}

	if err := scan.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}
