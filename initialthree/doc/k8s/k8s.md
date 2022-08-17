 k8s安装
 
## 安装教程

https://k8s.easydoc.net/docs/dRiQjyTY/28366845/6GiNOzyZ/9EX8Cp45

minikube
只是一个 K8S 集群模拟器，只有一个节点的集群，只为测试用，master 和 worker 都在一起
https://minikube.sigs.k8s.io/docs/start/



https://kubernetes.io/zh/docs/setup/production-environment/tools/kubeadm/install-kubeadm/

## 安装故障排查

```
kubeadm init --config=kubeadm-config.yaml
[kubelet-check] The HTTP call equal to 'curl -sSL http://localhost:10248/healthz' failed with error: Get "http://localhost:10248/healthz": dial tcp 127.0.0.1:10248: connect: connection refused.
```
https://zhuanlan.zhihu.com/p/406043822



https://kubernetes.io/zh/docs/setup/production-environment/tools/kubeadm/troubleshooting-kubeadm/

```
systemctl status kubelet
kubelet.service - kubelet: The Kubernetes Node Agent
   Loaded: loaded (/lib/systemd/system/kubelet.service; enabled; vendor preset: enabled)
  Drop-In: /etc/systemd/system/kubelet.service.d
           └─10-kubeadm.conf
   Active: activating (auto-restart) (Result: exit-code) since Mon 2022-03-21 11:13:35 CST; 1s ago
     Docs: https://kubernetes.io/docs/home/
  Process: 21601 ExecStart=/usr/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_CONFIG_ARGS $KUBELET_KUBEADM_ARGS $KUBELET_EXTRA_ARGS (code=exited, status=1/FAILURE)
 Main PID: 21601 (code=exited, status=1/FAILURE)

Mar 21 11:13:35 VM-3-10-ubuntu systemd[1]: kubelet.service: Unit entered failed state.
Mar 21 11:13:35 VM-3-10-ubuntu systemd[1]: kubelet.service: Failed with result 'exit-code'.
```
https://github.com/easzlab/kubeasz/issues/319

```
kubeadm init --config=kubeadm-config.yaml
[wait-control-plane] Waiting for the kubelet to boot up the control plane as static Pods from directory "/etc/kubernetes/manifests". This can take up to 4m0s
[kubelet-check] Initial timeout of 40s passed.
```
https://blog.csdn.net/weixin_44789466/article/details/119046245
https://q.cnblogs.com/q/124859/

```
systemctl status kubelet
kubelet.go:2001] "Skipping pod synchronization" err="PLEG is not healthy: pleg has yet to be successful"
```

获取pod参数

https://jimmysong.io/kubernetes-handbook/appendix/tricks.html

https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/

### 定义

Service 举个例子，考虑一个图片处理后端，它运行了 3 个副本。这些副本是可互换的 —— 前端不需要关心它们调用了哪个后端副本。 然而组成这一组后端程序的 Pod 实际上可能会发生变化， 
前端客户端不应该也没必要知道，而且也不需要跟踪这一组后端的状态。Service 定义的抽象能够解耦这种关联。
        



### 项目


使用k8s部署服务节点，每个服务单独运行在容器，解决了需要给每个服务分配启动端口的麻烦。不需要单独配置端口，由k8s分配一个内网地址。
集群中的端口不需要暴露出去，任通过center发现。对外的端口可以暴露出去，通过service把多个服务的端口映射到主机的一个端口。 
如果外部想链接指定的服务呢？


由于服务需要互联，还是需要告知本服务的IP。center 的IP需要用 service代理，固定下它的IP。多个center副本时是基于负载均衡，所以服务可能会
链接到不同的center，那它们就不能发现对方。 就只能用一个center？


配置文件中的端口可以固定下来，但是逻辑地址要有变化，要根据statefulSet来根据规则产生

