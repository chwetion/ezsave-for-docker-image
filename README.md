# EZSave
EZSave is a CLI tool that conversion between compressed packages and images.

### Getting Started
Download CLI and move it to $PATH
`chmod +x ezsave && mv ezsave /usr/local/bin`

### Example
#### Pull image and save
1. create your yaml file:
```yaml
packages:
# Specify the file name of the compressed package
- file: <compressed package name> # e.g. example.tar
  # The image contained in the compressed package
  content:
  - name: <re-tag name> # e.g. openjdk:7
    from: <image where pull> # e.g. <private-repository>/<project>/openjdk:7
  - name: <retag name> # e.g. openjdk:8
    from: <image where pull> # e.g. <private-repository>/<project>/openjdk:8
  - name: <retag name> # e.g. openjdk:11
    from: <image where pull> # e.g. <private-repository>/<project>/openjdk:11
```
2. ezsave save -f <filename>
3. you will get a example.tar file in execute directory

### Save image, re-tag, push
TBD