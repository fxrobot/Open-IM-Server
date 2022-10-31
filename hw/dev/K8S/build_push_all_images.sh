#/bin/sh
source ./path_info.cfg

# images version
version=v2.3.0
git pull
cd ../../script/; sh ./build_all_service.sh
cd ../hw/K8S/

for i in  ${service[*]}
do
  mv ../../bin/open_im_${i} ./${i}/
done

echo "move success"

echo "start to build images"

for i in ${service[*]}
do
	echo "start to build images" $i
	cd $i
	image="openim/${i}:$version"
	docker build -t $image . -f ./${i}.Dockerfile
	echo "build ${dockerfile} success"
    docker tag ${image} swr.cn-north-4.myhuaweicloud.com/wcmdh/${image}
	docker push swr.cn-north-4.myhuaweicloud.com/wcmdh/${image}
	echo "push ${image} success "
	cd ..
done

