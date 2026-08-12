## ONV環境構築メモ

必要要素のインストール
```bash
sudo apt update
sudo apt install ovn-central ovn-host openvswitch-switch
```
- ovn-central: ovn-northd + NB/SB DB(ovsdb-server)
- ovn-host: Node で動くエージェント
- openvswitch-switch: OVS

起動状態確認。
```bash
sudo systemctl status ovn-central
sudo systemctl status ovn-host
sudo systemctl status openvswitch-switch
```

多分DBソケットが `/var/run/ovn/ovnnb_db.sock` として待機しているはず。

統合Bridgeも自動で作成されているはずなので、後は設定した準備を
公式リファレンスに従って作成するのみ。

```bash
sudo ovs-vsctl set open . external-ids:ovn-remote="unix:/var/run/ovn/ovnsb_db.sock"
sudo ovs-vsctl set open . external-ids:ovn-encap-type=geneve
sudo ovs-vsctl set open . external-ids:ovn-encap-ip=127.0.0.1
sudo ovs-vsctl set open . external-ids:system-id=$(hostname)
```

Socketは、OVN界隈ではRootでしか触らないので、アプリから触る際にはEndpointが妥当。
```bash
sudo ovn-nbctl set-connection ptcp:6641:127.0.0.1
```
でTCP受付を行って、アプリエンドポイントもそこにする。

## 補足
当然、他のVM関連も入ってないといけない。
```bash
sudo apt update
sudo apt install libvirt-dev pkg-config build-essential
sudo apt install libvirt-daemon-system libvirt-clients qemu-system-x86 virtinst
```

Libvirtの管理グループに入っていないがゆえに起動できないことがある。
```bash
sudo usermod -aG libvirt $USER
```
即時反映はされないので、再ログイン必要。