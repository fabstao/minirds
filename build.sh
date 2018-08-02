# Do not run manually!

make all
docker build -t gcr.io/fabs-cl-02/mrds1:${BUILD_NUMBER} .
gcloud docker -- push gcr.io/fabs-cl-02/mrds1:${BUILD_NUMBER}
kubectl delete deployment mrds1
kubectl delete svc mrds1
kubectl run mrds1 --image=gcr.io/fabs-cl-02/mrds1:${BUILD_NUMBER} --port=8800 --replicas=1
kubectl expose deployment mrds1 --port=8800 --target-port=8800 --type=LoadBalancer
sleep 5
kubectl get deployments
sleep 60
kubectl get svc
