# EngOps deployment documentation

```
glide i
make test
make
```

```
docker build -t engops/ddvote .
docker tag engops/ddvote gcr.io/keen-autumn-144321/engops/ddvote
gcloud docker push gcr.io/keen-autumn-144321/engops/ddvote
```
