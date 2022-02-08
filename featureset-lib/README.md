# Mastro Featurestore client library

## Install

To install you can use pip, e.g.:

```
pip install mastro-fs==0.1.0
```
See [Pypi](https://pypi.org/project/mastro-fs/0.1.0/) for more details.

## Library Development

You can create a venv or a conda env respectively as follows:

```bash
virtualenv mastroenv
python -m venv mastroenv
source mastroenv/bin/activate
```

```bash
conda create -n mastroenv python=3.9
conda activate mastroenv
```

You can then start use the local env and develop with the local version of the library.

## Examples

The examples were prepared using Jupyter with a local version of the library. 
Please check the install section on how to install the library.

To run a local version, you can add a new environment along with a new jupyter kernel:
```bash
python -m pip install ipykernel
python -m ipykernel install --user --name mastroenv --display-name "Python (mastroenv)
```

Please have a look at [the example notebook](examples/mastrofs.ipynb).
