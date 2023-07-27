# commands used for set up

## cobra-cli installation

```bash
go install github.com/spf13/cobra-cli@latest
```

## project set up

```bash
go mod init futurama
cobra-cli init futurama --author "Aric Hansen <aric.p.hansen@gmail.com>" --license mit
```

## adding commands

```bash
cd futurama
cobra-cli add get
cobra-cli add quote -p 'getCmd'
```

## futurama installation

```bash
go install github.com/aric-h/futurama@latest
```
