from setuptools import find_packages, setup

setup(
    name='mastro_fs',
    description='A Python library to interact with Mastro feature stores',
    keywords = ['mastro', 'featurestore', 'mlops'],
    packages=find_packages(include=['mastro_fs']),
    version='0.1.0',
    author='pilillo',
    license='Apache License 2.0',
    install_requires=[],
    setup_requires=['pytest-runner==5.3.1'],
    tests_require=['pytest==6.2.5'],
    test_suite='tests'
)