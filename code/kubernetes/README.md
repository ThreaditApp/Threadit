```sh
gcloud config set project fcul2025
PROJECT_ID=$(gcloud config get-value project)
IMAGE_TAG=latest

gcloud container clusters create threadit-cluster --num-nodes=2 --machine-type=e2-standard-2 --region=europe-west1 --disk-type=pd-standard --disk-size=100

gcloud container clusters get-credentials threadit-cluster --region=europe-west1

gcloud auth configure-docker

cd code
docker build -t gcr.io/$PROJECT_ID/db-service:$IMAGE_TAG -f services/db-service/Dockerfile .
docker build -t gcr.io/$PROJECT_ID/community-service:$IMAGE_TAG -f services/community-service/Dockerfile .
docker build -t gcr.io/$PROJECT_ID/grpc-gateway:$IMAGE_TAG -f grpc-gateway/Dockerfile .

docker push gcr.io/$PROJECT_ID/db-service:$IMAGE_TAG
docker push gcr.io/$PROJECT_ID/community-service:$IMAGE_TAG
docker push gcr.io/$PROJECT_ID/grpc-gateway:$IMAGE_TAG

gcloud compute disks create dataset-disk --size=5GB --zone=europe-west1-b
gcloud compute disks create mongodb-disk --size=10GB --zone=europe-west1-b

cd kubernetes

kubectl apply -f configmap.yml
kubectl apply -f secrets.yml
kubectl apply -f traefik-configmap.yml

kubectl apply -f dataset-pv.yml
kubectl apply -f dataset-pvc.yml

kubectl apply -f copy-dataset-pod.yml
kubectl cp ../dataset copy-dataset-pod:/dataset
kubectl exec -it copy-dataset-pod -- ls /dataset
kubectl exec -it copy-dataset-pod -- cp /dataset/dataset/communities.json /dataset/communities.json
kubectl exec -it copy-dataset-pod -- cp /dataset/dataset/threads.json /dataset/threads.json
kubectl delete pod copy-dataset-pod

kubectl apply -f mongodb-service.yml
kubectl apply -f db-service-service.yml
kubectl apply -f community-service-service.yml
kubectl apply -f grpc-gateway-service.yml
kubectl apply -f traefik-service.yml

kubectl apply -f mongodb-statefulset.yml
kubectl apply -f db-service-deployment.yml
kubectl apply -f community-service-deployment.yml
kubectl apply -f grpc-gateway-deployment.yml
kubectl apply -f traefik-deployment.yml

# Grab external ip and test request (e.g. http://34.78.232.130/api/communities)
kubectl get service traefik

# Example: View all resources
kubectl get all

# Example: Check pod logs
kubectl logs -f deployment/community-service
```
