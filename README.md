## IN PROGRESS

### Install
Inspiration from https://github.com/vfarcic/metacontroller-demo
Metacontroller taken from https://github.com/metacontroller/metacontroller/tree/master/manifests/production


kubectl apply -f metacontroller-crds-vi.yaml 
kubectl apply -f metacontroller-namespace.yaml 
kubectl apply -f metacontroller-rbac.yaml  
kubectl apply -f metacontroller.yaml 

kubectl apply -f crds.yaml
kubectl apply -f composite-controllers.yaml
kubectl apply -f controllers.yaml


