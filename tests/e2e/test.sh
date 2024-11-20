#!/bin/bash

INPUT=${1}
EXPECT_STATUS=${2}

COMMIT_HASH=$(git show --format='%H' --no-patch)
docker image build -t e2e:$COMMIT_HASH .

kubectl apply -f cron.yaml
kubectl patch -f cron.yaml -p '{"spec":{"jobTemplate":{"spec":{"template":{"spec":{"containers":[{"name":"stub-cronjob","image":"e2e:'$COMMIT_HASH'"}]}}}}}}'

kubectl create transient_configmap stub-configmap --from-literal=num=$INPUT --job-name=test --job-from=cronjob/stub-cronjob | tee /dev/null
if [ ${PIPESTATUS[0]} -ne $EXPECT_STATUS ]; then
    echo "unexpected status"
    exit 1
fi
kubectl get configmap stub-configmap | /dev/null
if [ ${PIPESTATUS[0]} -ne 1 ]; then
    echo "configmap has not been deleted"
    exit 1
fi
kubectl get job test | /dev/null
if [ ${PIPESTATUS[0]} -ne 1 ]; then
    echo "job has not been deleted"
    exit 1
fi

echo "test passed"