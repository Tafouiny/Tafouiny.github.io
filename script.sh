#Create Image from dockerfile
docker build -t forum .

#List Images
#docker images

#Create Container from forum image
docker run --name cont -d -p 8080:8080 forum

#List Containers
#docker ps -a

#Stop container and remove
#docker stop cont
#docker rm cont

#Delete image
#docker rmi forum