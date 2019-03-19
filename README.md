# BOM Resource

Checks out a specific commit from a git `product` repository and optionally adds features file into `source` folder

## Source Configuration

* `uri`: *Required.* The location of the repository.

* `branch`: The branch to track. This is *optional* if the resource is
   only used in `get` steps (default value in this case is `master`).
   However, it is *required* when used in a `put` step.

* `private_key`: *Optional.* Private key to use when pulling/pushing.
    Example:
    ```
    private_key: |
      -----BEGIN RSA PRIVATE KEY-----
      MIIEowIBAAKCAQEAtCS10/f7W7lkQaSgD/mVeaSOvSF9ql4hf/zfMwfVGgHWjj+W
      <Lots more text>
      DWiJL+OFeg9kawcUL6hQ8JeXPhlImG6RTUffma9+iGQyyBMCGd1l
      -----END RSA PRIVATE KEY-----
    ```
    Note: You can also use pipeline templating to hide this private key in source control. (For more information: https://concourse.ci/fly-set-pipeline.html)

* `username`: *Optional.* Username for HTTP(S) auth when pulling/pushing.
  This is needed when only HTTP/HTTPS protocol for git is available (which does not support private key auth)
  and auth is required.

* `password`: *Optional.* Password for HTTP(S) auth when pulling/pushing.

   Note: You can also use pipeline templating to hide this password in source control. (For more information: https://concourse.ci/fly-set-pipeline.html)

* `skip_ssl_verification`: *Optional.* Skips git ssl verification by exporting
  `GIT_SSL_NO_VERIFY=true`.

* `git_config`: *Optional.* If specified as (list of pairs `name` and `value`)
  it will configure git global options, setting each name with each value.

  This can be useful to set options like `credential.helper` or similar.

  See the [`git-config(1)` manual page](https://www.kernel.org/pub/software/scm/git/docs/git-config.html)
  for more information and documentation of existing git options.

* `commit_verification_keys`: *Optional.* Array of GPG public keys that the
  resource will check against to verify the commit (details below).

* `commit_verification_key_ids`: *Optional.* Array of GPG public key ids that
  the resource will check against to verify the commit (details below). The
  corresponding keys will be fetched from the key server specified in
  `gpg_keyserver`. The ids can be short id, long id or fingerprint.

* `gpg_keyserver`: *Optional.* GPG keyserver to download the public keys from.
  Defaults to `hkp:///keys.gnupg.net/`.

* `https_tunnel`: *Optional.* Information about an HTTPS proxy that will be used to tunnel SSH-based git commands over.
  Has the following sub-properties:
    * `proxy_host`: *Required.* The host name or IP of the proxy server
    * `proxy_port`: *Required.* The proxy server's listening port
    * `proxy_user`: *Optional.* If the proxy requires authentication, use this username
    * `proxy_password`: *Optional.* If the proxy requires authenticat, use this password

* `bom_root_directory`: *Optional.* Path to use within bom repository

* `features`: *Optional.* Allows configuration of how feature files will be compiled.
  Has the following sub-properties:
  * `validate_keys`: *Optional.* Whether keys should be validated based on expected_keys file in product repo, default: `false`
  * `expected_keys_file`: *Optional.* Specifies name of expected keys yaml file.  default: `expectedKeys.yml`
  * `directory`: *Optional.* Directory within product repo that contains feature files.  default: `features`
  * `search_path`: *Optional.* Array of features files in `directory` that will be added to search path

### Example

Resource configuration for a private repo with an HTTPS proxy:

``` yaml
resource_types:
- name: bom
  type: docker-image
  source:
    repository: pivotalservices/bom-resource
    tag: latest

resources:
- name: source-code
  type: bom
  source:
    uri: git@github.com:pivotalservices/bom-resource.git
    branch: master
    private_key: |
      -----BEGIN RSA PRIVATE KEY-----
      MIIEowIBAAKCAQEAtCS10/f7W7lkQaSgD/mVeaSOvSF9ql4hf/zfMwfVGgHWjj+W
      <Lots more text>
      DWiJL+OFeg9kawcUL6hQ8JeXPhlImG6RTUffma9+iGQyyBMCGd1l
      -----END RSA PRIVATE KEY-----
    git_config:
    - name: core.bigFileThreshold
      value: 10m
    disable_ci_skip: true
    https_tunnel:
      proxy_host: proxy-server.mycorp.com
      proxy_port: 3128
      proxy_user: myuser
      proxy_password: myverysecurepassword
    features:
      validate_keys: false
      directory: features
      expected_keys_file: foo.yml
      search_path: ["global.yml", "us.yml", "production.yml"]
```

Cloning cf product with a cf-vars.yml feature file:

``` yaml
- get: source-code
  params:
    product: cf
```

## Behavior

### `check`: Check for new commits.

The repository is cloned (or pulled if already present), and any commits
from the given version on are returned. If no version is given, the ref
for `HEAD` is returned.

### `in`: Clone the repository, at the given ref and if product name is specified will clone that product in it's place.

Clones the repository to the destination, and locks it down to a given ref.
It will return the same given ref as version.

Submodules are initialized and updated recursively.

#### Parameters

* `depth`: *Optional.* If a positive integer is given, *shallow* clone the
  repository using the `--depth` option. Using this flag voids your warranty.
  Some things will stop working unless we have the entire history.

* `submodules`: *Optional.* If `none`, submodules will not be
  fetched. If specified as a list of paths, only the given paths will be
  fetched. If not specified, or if `all` is explicitly specified, all
  submodules are fetched.

* `disable_git_lfs`: *Optional.* If `true`, will not fetch Git LFS files.

* `product`: *Optional.* If set, will clone product repository based on configuration of product bom file.

#### GPG signature verification

If `commit_verification_keys` or `commit_verification_key_ids` is specified in
the source configuration, it will additionally verify that the resulting commit
has been GPG signed by one of the specified keys. It will error if this is not
the case.

#### Additional files populated

 * `.git/committer`: For committer notification on failed builds.
   This special file `.git/committer` which is populated with the email address
   of the author of the last commit. This can be used together with  an email
   resource like [mdomke/concourse-email-resource](https://github.com/mdomke/concourse-email-resource)
   to notify the committer in an on_failure step.

 * `.git/ref`: Version reference detected and checked out. It will usually contain
   the commit SHA-1 ref, but also the detected tag name when using `tag_filter`.

 * `.git/bom`: the source of the bom used when cloning the product repo.

 * `.git/features`: compiled feature file with yaml content

## BOM Specifics

### Repository configuration
To enable the bom resource to use a single concourse resource that can clone different git repositories need to configure a repository with the following format.

```
├── bom
│   └── test.yml
├── features
│   └── test-vars.yml
```

Where test represents a `product` that can be specified as a params on the `get` of the resource.

### BOM file format
The bom file needs to have at mininum 2 properties

* `git-repo`: Uri to the git repo for this given product.  Support both https and ssh format uri.
* `commit`: Represent the commit to checkout in the `git-repo` specified.  

#### Example

``` yaml
commit: eb90453b0ca5768fa6....
git-repo: https://github.com/pivotalservices/bom-resource.git
```

### Feature File
Feature files are supported as a simple flat structure of key/values in a yaml file

#### Example
``` yaml
foo: bar
hello: world
```

#### Search order
Feature files will be searched for and applied in the following order.  This uses a last-in wins key model.
- Within the product repository files specified in the `search_path` property in the `directory` folder based on configuration of the `features` within `source` configuration of the resource
- Within the bom repository `features/global-vars.yml`
- Within the bom repository `features/<product>-vars.yml`

This allows to setup property inheritance model for feature values that can be used within your concourse tasks.

#### Expected Keys
To have keys validated need to include a yaml file in the "product repository" with the following format.

``` yaml
keys:
- foo
- hello
```

## Development

### Prerequisites

* golang is *required* - version 1.9.x is tested; earlier versions may also
  work.
* docker is *required* - version 17.06.x is tested; earlier versions may also
  work.

### Running the tests

The tests have been embedded with the `Dockerfile`; ensuring that the testing
environment is consistent across any `docker` enabled platform. When the docker
image builds, the test are run inside the docker container, on failure they
will stop the build.

Run the tests with the following command:

```sh
docker build -t pivotalservices/bom-resource .
```

### Contributing

Please make all pull requests to the `master` branch and ensure tests pass
locally.
