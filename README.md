## images.social

Image uploading service. Nothing fancy involved.

#### TODO

- gossip between instances and can transfer files between each other for serving
  - otherwise, redirect over to who has upload
  - Immediate serving is handled by the server which got an upload responds, other reqs can fail until transfer finishes
    - Submit reqs to other instances to transfer files over there
