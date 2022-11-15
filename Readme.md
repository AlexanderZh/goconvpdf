# Simple restful service for document conversion

Uses libreoffice to convert from any document (txt, doc, docx, xls, xlsx, rtf, jpg, ...) to pdf

Dirty and ugly wrapper

## Usage with docker:

1. Run service

```bash
docker build -t goconvpdf .
docker run -p 8080:8080 --rm --mount type=tmpfs,destination=/tmpfs goconvpdf
```

Using  ```--mount type=tmpfs,destination=/tmpfs``` is optional: it provides ramdisk instead of HDD and might be slightly faster for large files.

2. Send your file via curl or custom HTTP call:

```bash
curl -XPOST localhost:8080 -H "Content-Type: multipart/form-data" -F fileName=@routers.go --output res.pdf
```

res.pdf - is your converted file