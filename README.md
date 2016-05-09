
# Travix Fireball ADK 
App Developer Kit for the Travix Fireball infrastructure. The ADK consists of Appix, a CLI tool to publish Apps to the App Catalog. 

## Docker
For your convenience, a Dockerfile is created takes care of the integration into the NH website for you and also includes the Yeoman generator to scaffold new Apps. 

Create the Docker image using `docker build -t adk .`. Make sure that your Docker machine has enough space available, the entire NH repository tree is retrieved, including content pages.

Run the Docker container using `docker run -v /path/to/your/apps/folder:/c/home/docker/apps -it adk bash`. The NH website will be started when you run the docker container. Usage instructions are displayed at startup.

## Appix CLI
After building (bin/build.bat / bin/build.sh), run `bin/appix --help` to view the commands and usage.
