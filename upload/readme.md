#JCLI upload

This program is what you'd use to select a manuscript file from your local machine and upload it to a submission system. In a fully realised system this would involve uploading it to some kind of cloud storage, but at this stage we'll make it do very little.

*Given* That I have a manuscript file on my computer
*And* The file extension of that is '.md'
*When* I supply the file as an argument to `jcli-upload`
*Then* It should echo the file/path so that it could be read by a successor program

*Given* That I have a manuscript file on my computer
*And* The file extension is not '.md'
*When* I supply the file as an argument to `jcli-upload`
*Then* It should exit with an error and explain that at present only markdown files are supported
