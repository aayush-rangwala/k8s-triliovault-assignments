# Assignment-4

### Reconciling on advanced set of restrictions:

1. To run the controller in one namespace, I added `namespace:` field in `main.go` like this:

`namespace:="vishu"`

``` 
mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		Namespace:          namespace,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
	})
```

2. To make a 3 level ownership (with 3 level resources), I first made a new function `buildDeploymentPodJob()` that will create a deployment, then this deployment will create a pod and that pod will create a job from a script command provided to it. Since ownership was required on one source only, I added the `OwnerReferences` in deployment only.

* To make a Kubernetes Job from a pod, I made `job.go` file which takes the `kubeconfig` as in-cluster and creates, list and deletes a job. I built the binary of file as `job` and made a Dockerfile which contains the binary. This docker file is converted to docker image and pushed to docker hub. Now, this image is referred to the actual container image of the pod that is going to run that job.

Steps performed to make the Kubernetes Job, make its binary, build the docker image and use it as the our main pod image:

- `kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default`
- Made job.go file
- Made Dockerfile
- Installed all the dependencies that were required with job.go file (like client-go, googleapis-gnostic)
- GOOS=linux go build -o ./app . (exporting the env variable and building go code)
- docker build -t jobimage . (building dockerfile for pulling it afterwards in pods)
- Pushed the docker image to docker hub. (using docker login)
- Used the docker image name as our main pod image. (docker.io/username/imagename)

3. Did CRUD operations on all the 3 levels of resources (deployment, pod, job).

4. Used `mgr.GetFieldIndexer().IndexField()` with all the resource whether they had `OwnerReferences` or not. 
* Checked for `API Version` on deployment.
* Checked for `API Version` & `Kind` of Pod and Job.

5. To check that the resource are working inside the given namespace, I used a predicate that will filter and check for the resource running the given namespace like this:

```
checkNamespace := func(e event.CreateEvent) bool {
		if e.Meta.GetNamespace() != "vishu" {
			log.Fatal("resource not running in vishu namespace")
			return false
		}
		return true
	}
```

```	
p := predicate.Funcs{
		CreateFunc: checkNamespace,
	} 
```

` WithEventFilter(p) `