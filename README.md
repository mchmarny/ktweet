# ktweet

Knative Twitter Source

Simple Twitter event source implementation for [Knative Eventing](http://github.com/knative/eventing) using `ContainerSource`.

## Setup

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
For example, to have the `ktweet` produce events for matching the Twitter search for the term `Knative` your
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
kubectl apply -f source.yaml -n demo
```

### Logs

Once you know there are some tweets matching your search you might have to wait few seconds for `ktweet`
fetch it and you should see it in your target service

```shell
kubectl logs -l eventing.knative.dev/source=twitter-source -n demo -c source
```
