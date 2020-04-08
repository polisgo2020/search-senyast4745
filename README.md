# Techopolis Golang homework
Golang Open Source Homework Repository

## Build application

```shell script
make build
```

## Usage

### Run build index

```shell script
./search build --soruces /path/to/folder/to/index --index /index/file/path
```

### Run search

The program can be launched in two ways

#### Console search

```shell script
export LISTEN=inteface-to-listen
export LOG_LEVEL=log-level
export TIMEOUT=server-timeout  
./search search --index /index/file/path
```

After it you can send request to server
```http request
POST /?search=`search-phrase` HTTP/1.1
Host: `interfase-to-listen`
```

#### Search in docker

Before all run **Elastic Stack**:

```shell script
docker-compose up -f --build docker-compose-logs.yml
```

After you can up search docker-compose:

```shell script
export PORT=8080
export LOG_LEVEL=`log-level`
export IND_FOLDER=/path/to/folder/with/index-file
export IND_FILE=/index/file/name
docker-compose up --build
```
You can go in your browser to `localhost` and start searching by web-interface 