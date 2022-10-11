#!/bin/sh
set -e
if [ -f $WORKSPACE/../TOGGLE ]; then
    echo "****************************************************"
    echo "data.stack.b2b.agent.watcher :: Toggle mode is on, terminating build"
    echo "data.stack.b2b.agent.watcher :: BUILD CANCLED"
    echo "****************************************************"
    exit 0
fi

cd $WORKSPACE

cDate=`date +%Y.%m.%d.%H.%M` #Current date and time

if [ -f $WORKSPACE/../DATA_STACK_RELEASE ]; then
    REL=`cat $WORKSPACE/../DATA_STACK_RELEASE`
fi
BRANCH='dev'
if [ -f $WORKSPACE/../BRANCH ]; then
    BRANCH=`cat $WORKSPACE/../BRANCH`
fi
if [ $1 ]; then
    REL=$1
fi
if [ ! $REL ]; then
    echo "****************************************************"
    echo "data.stack.b2b.agent.watcher :: Please Create file DATA_STACK_RELEASE with the releaese at $WORKSPACE or provide it as 1st argument of this script."
    echo "data.stack.b2b.agent.watcher :: BUILD FAILED"
    echo "****************************************************"
    exit 0
fi
TAG=$REL

echo "****************************************************"
echo "data.stack.b2b.agent.watcher :: Using build :: "$TAG
echo "****************************************************"

echo "****************************************************"
echo "data.stack.b2b.agent.watcher :: Adding IMAGE_TAG in Dockerfile :: "$TAG
echo "****************************************************"
sed -i.bak s#__image_tag__#$TAG# Dockerfile
sed -i.bak s#__signing_key_user__#$SIGNING_KEY_USER# Dockerfile
sed -i.bak s#__signing_key_password__#$SIGNING_KEY_PASSWORD# Dockerfile

if [ -f $WORKSPACE/../CLEAN_BUILD_B2B_AGENT_WATCHER ]; then
    echo "****************************************************"
    echo "data.stack.b2b.agent.watcher :: Doing a clean build"
    echo "****************************************************"

    docker build --no-cache -t data.stack.b2b.agent.watcher.$TAG .
    rm $WORKSPACE/../CLEAN_BUILD_B2B_AGENT_WATCHER

else
    echo "****************************************************"
    echo "data.stack.b2b.agent.watcher :: Doing a normal build"   
    echo "****************************************************"
    docker build -t data.stack.b2b.agent.watcher.$TAG .
fi
echo "****************************************************"
echo "data.stack.b2b.agent.watcher :: BUILD SUCCESS :: data.stack.b2b.agent.watcher.$TAG"
echo "****************************************************"
echo $TAG > $WORKSPACE/../LATEST_B2B_AGENT_WATCHER
