# t7s

Minimal templating engine. Create your meta configuration files with ease. Supports indexing of templates to generate a variables file, enabling you to build your templates first and then generate and fill in the variables later.

## Usage

```bash
Usage: t7s [options] [in] [out]
Options:
  -i    Index templates
  -left string
        Left delimiter (default "{{")
  -log string
        Log level (default "info")
  -m    Index merge with existing variables
  -r    Require all variables to be set
  -right string
        Right delimiter (default "}}")
  -v string
        Variables file (default "variables.yaml")
  -version
        Print version
```

## Variables File

The variables file is a YAML file that contains a list of variables that will be used to render the template. The variables file can be generated by indexing the template file(s) or created manually.

```yaml
variables:
- name: name
  value: world
  description: name of the person to greet
  required: true
```

The `description` and `required` fields are optional. If `required` is set to `true`, then the variable must be set in the variables file or an error will be thrown, regardless of whether the `-r` flag is set.

Variable values can be set literally in the variables file or can reference an environment variable using either `$NAME` or `${NAME}` syntax.

```yaml
variables:
- name: name
  value: $NAME
  description: name of the person to greet
```

## Examples

### Simple

```bash
# create a template file
$ cat > template.txt <<EOF
Hello {{name}}!
EOF

# create a variables file
$ cat > variables.yaml <<EOF
variables:
- name: name
  value: world
  description: name of the person to greet
EOF

# render the template
$ t7s template.txt
Hello world!
```

### Index

```bash
# create a template file
$ cat > template.txt <<EOF
Hello {{name}}!
EOF

# index the file to generate a variables file
$ t7s -i template.txt variables.yaml

# edit the variables file
$ cat variables.yaml
variables:
- name: name
  value: world

# render the template
$ t7s template.txt
Hello world!
```

### Index Merge

```bash
# create a template file
$ cat > template.txt <<EOF
Hello {{name}}!
This is my unindexed variable: {{address}}
EOF

# create a variables file
$ cat > variables.yaml <<EOF
variables:
- name: name
  value: world
  description: name of the person to greet
EOF

# index the file to generate a variables file, merging with the existing variables
$ t7s -i -m template.txt variables.yaml

# cat the variables file
$ cat variables.yaml
variables:
- name: name
  value: world
  description: name of the person to greet
- name: address
  value: ""
```

### Require

```bash
# create a template file
$ cat > template.txt <<EOF
Hello {{name}}!
EOF

# create a variables file, missing the required variable
$ cat > variables.yaml <<EOF
variables:
- name: name
  description: name of the person to greet
EOF

# render the template
$ t7s -r template.txt
Error: variable 'name' is not set
```

### Delimiters

```bash
# create a template file, using different delimiters
$ cat > template.txt <<EOF
Hello ||name||!
EOF

# create a variables file
$ cat > variables.yaml <<EOF
variables:
- name: name
  value: world
  description: name of the person to greet
EOF

# render the template, using different delimiters
$ t7s -left "||" -right "||" template.txt
Hello world!
```

### Directory Support

If a directory is supplied as the input, it will automatically be recursively traversed and all files will be rendered/indexed.

```bash
# create a template directory
$ mkdir templates

# create a template file
$ cat > templates/template.txt <<EOF
Hello {{name}}!
EOF

# create another template file
$ cat > templates/another-template.txt <<EOF
Hello {{another-name}}!
EOF

# create a variables file
$ cat > variables.yaml <<EOF
variables:
- name: name
  value: world
  description: name of the person to greet
- name: another-name
  value: world
  description: name of the person to greet
EOF

# render the template directory
$ t7s templates /tmp/generated
$ cat /tmp/generated/template.txt
Hello world!
$ cat /tmp/generated/another-template.txt
Hello world!
```

### Indexing a Directory

If a directory is supplied as the input, it will automatically be recursively traversed and all files will be rendered/indexed.

```bash
# create a template directory
$ mkdir templates

# create a template file
$ cat > templates/template.txt <<EOF
Hello {{name}}!
EOF

# create another template file
$ cat > templates/another-template.txt <<EOF
Hello {{another-name}}!
EOF

# index the template directory to generate a variables file
$ t7s -i templates variables.yaml
```