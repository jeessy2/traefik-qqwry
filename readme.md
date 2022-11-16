- download qqwry.dat

  ```shell
  curl -o qqwry.dat https://99wry.cf/qqwry.dat
  ```

- config

  ```yaml
    # Static configuration
    experimental:
      plugins:
        example:
          moduleName: github.com/jeessy2/traefik-qqwry
          version: v0.0.1
  ```

  ```yaml
  http:
    middlewares:
      my-plugins:
        plugin:
          traefik-qqwry:
            dbPath: /opt/plugins-storage/qqwry.dat
            # headers:
            #  city: "X-City"
            #  isp: "X-Isp"
  ```

- thanks
  - https://github.com/coderjanus/encoding
  - https://github.com/xiaoqidun/qqwry