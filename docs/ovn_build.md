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