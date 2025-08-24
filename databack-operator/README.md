# databack-operator
数据备份，基于k8s crd+operator实现

环境说明：
- 开发环境：mac go1.18.x
- 运行环境：linux amd64(k8s master--kubectl环境，操作k8s)

## 项目初始化

基于kubebuilder完成初始化（kubebuilder是一个crd+operator脚手架）
```bash
https://github.com/kubernetes-sigs/kubebuilder.git
```
编译kubebuilder
```bash
make build
```
基于kubebuilder初始化代码结构
```bash
kubebuilder init databack-operator --domain="operator.kubeimooc.com" --project-name="databack-operator" --repo="kubeimooc.com/operator-databackup"
```
创建API

``` bash
kubebuilder create api --group "" --version v1beta1 --kind Databack
```

## 项目发布

安装crd
```bash
make install
```
打包operator
```bash
make docker-build docker-push IMG=harbor.kubeimooc.com/operator/databack:v1beta1
```
发布operator
```bash
make deploy IMG=harbor.kubeimooc.com/operator/databack:v1beta1
```
部署databack服务
```bash
kubectl apply -f config/samples/_v1beta1_databack.yaml
```

## 项目开发

### 定义crd
修改：/api/v1beta1/databack_types.go

### 开发operator
修改：/controllers/databack_controller.go
实现Reconcile的逻辑

## 正式对外提供
导出对外提供服务的资源文件
```bash
mkdir pushlish
/root/databack-operator/bin/kustomize build ../config/default > databack-operator.yaml
```