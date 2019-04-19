# ktweet

Knative Twitter Source

Simple Twitter event source implementation for [Knative Eventing](http://github.com/knative/eventing) using `ContainerSource`.

## Setup

### Istio Mesh

> NOte, Until Istio 1.1 you'll need to annotate the
> twitter source with the `traffic.sidecar.istio.io/includeOutboundIPRanges`
> to ensure the source can emit events into the mesh.

To configure [outbound network access](https://github.com/knative/docs/blob/master/docs/serving/outbound-network-access.md) you will need to determine the IP range of your cluster and capture it into
your `NET_SCOPE` variable.

```shell
export NET_SCOPE=$(gcloud container clusters describe ${CLUSTER_NAME} --zone=${CLUSTER_ZONE} \
                  | grep -e clusterIpv4Cidr -e servicesIpv4Cidr \
                  | sed -e "s/clusterIpv4Cidr://" -e "s/servicesIpv4Cidr://" \
                  | xargs echo | sed -e "s/ /,/")
```

End enter it into the  `config/source.yaml` file:

```yaml
kind: ContainerSource
metadata:
  annotations:
    traffic.sidecar.istio.io/includeOutboundIPRanges: "NET_SCOPE"
...
```

### Twitter API

To configure this event source you will need Twitter API access keys. [Good instructions on how to get them](https://iag.me/socialmedia/how-to-create-a-twitter-app-in-8-easy-steps/)



### Knative Secret

Create a secret with your Twitter access keys

```shell
kubectl create secret generic ktweet-secrets -n demo \
    --from-literal=T_CONSUMER_KEY=${T_CONSUMER_KEY} \
    --from-literal=T_CONSUMER_SECRET=${T_CONSUMER_SECRET} \
    --from-literal=T_ACCESS_TOKEN=${T_ACCESS_TOKEN} \
    --from-literal=T_ACCESS_SECRET=${T_ACCESS_SECRET}
```

## Run

To launch `ktweet` event source just define the Twitter search query in the `config/source.yaml`.
For example, to have the `ktweets` produce events for matching the Twitter search for the term `Knative` your
`config/source.yaml` would look like this:

```yaml
...
spec:
  args:
  - --query=Knative
  ...
```

When done editing, just apply into your Knative cluster:

```shell
kubectl apply -f config/source.yaml -n demo
```

### Logs

Once you know there are some tweets matching your search you might have to wait few seconds for `ktweets`
fetch it and you should see it in your target service

```shell
kubectl -l 'serving.knative.dev/service=twitter-viewer' logs -c user-container
```
