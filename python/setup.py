from logging import getLogger
import shutil
import subprocess

from setuptools import setup, find_packages
from setuptools.command.develop import develop
from setuptools.command.install import install

logger = getLogger(__name__)


def build_go_archive():
    logger.info("Building C archive with static library")
    if shutil.which("go") is None:
        raise EnvironmentError("You should have Go installed and available on your path in order to build this module")
    subprocess.check_output(["make", "static"], cwd="../")
    logger.info("C archive successfully built")


class build_static_and_develop(develop):
    
    def run(self):
        build_go_archive()
        super(build_static_and_develop, self).run()


class build_static_and_install(install):
    
    def run(self):
        build_go_archive()
        super(build_static_and_install, self).run()


setup(
    name="enry",
    version="0.1.1",
    description="Python bindings for go-enry package",
    setup_requires=["cffi>=1.0.0"],
    cffi_modules=["build_enry.py:ffibuilder"],
    packages=find_packages(),
    install_requires=["cffi>=1.0.0"],
    cmdclass={"develop": build_static_and_develop, "install": build_static_and_install}
)
