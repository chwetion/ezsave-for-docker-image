# EZSave
EZSave is a CLI tool that conversion between compressed packages and images.

### Getting Started
Download CLI and move it to $PATH
`chmod +x ezsave && mv ezsave /usr/local/bin`

### Example
#### Yaml file
```yaml
defaultTarget:
  registry: docker.io # default value is docker.io
  project: library # default value is library
defaultFrom:
  registry: docker.io # default value is docker.io
  project: library # default value is library
auths: # it will be used in push and pull when same registry address
- address: <registry address>
  username: <username>
  password: <password>
packages:
# Specify the file name of the compressed package
- file: <compressed package name> # e.g. example.tar
  # The image contained in the compressed package
  content:
  - name: <re-tag name> # e.g. openjdk:7
    from:
      registry: <private registry when pull> # e.g. x.x.x.x:port, will override default
      project: <project name when pull> # e.g. library, will override default
      name: <image name include tag when pull> # e.g. openjdk:7, will override default
    target:
      registry: <private registry when push> # e.g. y.y.y.y:port, will override default
      project: <project name when push> # e.g. base, will override default
      name: <image name include tag when push> # e.g. openjdk:7, will override default
  - name: <retag name> # e.g. openjdk:8, if not from field and target filed, will use defaultTarget and defaultFrom field fill, and name will use content[*].name
  - name: <retag name> # e.g. openjdk:11
```
#### Pull image and save
1. create your yaml file;
2. execute `ezsave save -f <filename>`;
3. you will get a example.tar file in execute directory. Image will re-tag be content.name field.

#### Save image, re-tag, push
1. create your yaml file;
2. execute `ezsave load -f <filename>`;
3. ezsave will load your compressed package which you specify, re-tag images, push target registry.

#### Pull image, re-tag, push
1. create your yaml file;
2. execute `ezsave forward -f <filename>`;
3. ezsave will pull images from list, re-tag, and push target registry.