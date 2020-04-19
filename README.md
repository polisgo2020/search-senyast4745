# invindex
github.com/polisgo2020/search-senyast4745 implements inverted index to perform full-text search.


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

#### Search + building index in docker

You can up invindex in docker-compose:

```shell script
mkdir ./data
mkdir -p ./mongo/data
cp /data/files/folder ./data
export LOG_LEVEL=<log-level>;DATABASE=<your-database-name>
export DB_USERNAME=<your-database-user>;DB_PASSWORD=<your-database-password>
export DB_INTERFACE=mongo://<DB_USERNAME>:<DB_PASSWORD>@<your-database-host>
docker-compose up -d
```

> To start the build and search in Docker, you must have the **./data** folder in the same directory as the **docker-compose.yml** file.
> Files to index should be in this (**./data**) folder.

After it you can go in your browser to [localhost](http://localhost) and start searching by web-interface.

#### Add Kibana logs

If you want to use [**Kibana**](https://www.elastic.co/kibana) to view application logs:

* Add code to `backend`, `building` and `nginx` service in ``docker-compose.yml``:
```yaml
    logging:
      driver: "fluentd"
      options:
        fluentd-address: localhost:24224
        tag: backend.log
```

* After run **Elastic Stack**:

```shell script
docker-compose -f docker-compose-logs.yml up -d
```

* Then repeat the steps from the stage: **"Search in docker"**.

> To secure your logs [see this guide.](http://codingfundas.com/setting-up-elasticsearch-6-8-with-kibana-and-x-pack-security-enabled/index.html)