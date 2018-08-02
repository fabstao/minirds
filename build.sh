# Do not run manually!

make all
docker build -t gcr.io/fabs-cl-02/mrds1:${BUILD_NUMBER} .
gcloud docker -- push gcr.io/fabs-cl-02/mrds1:${BUILD_NUMBER}
kubectl delete deployment mrds1
kubectl run mrds1 --image=gcr.io/fabs-cl-02/mrds1:${BUILD_NUMBER} --port=8800 --replicas=1
sleep 5
kubectl gte deployments
kubectl get svc
