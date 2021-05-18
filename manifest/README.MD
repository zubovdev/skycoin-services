## manifest Tool to store files' metadata

manifest is a CLI tool storing information about files(metadata) of a directory, to use the tool:
- execute 'manifest init' command in a directory to create the .cxo folder 
- then execute 'manifest commit' command to store the files metadata in a .cxo file in .cxo/checkpoints/ folder
- the commit command has the -print-json flag to also print the files info for you to review, and a flag -meta to include a meatadata section in json.
