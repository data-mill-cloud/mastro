clean:
	rm -rf dist/*

build: clean
	python setup.py bdist_wheel

upload: build
	python -m twine upload dist/*
