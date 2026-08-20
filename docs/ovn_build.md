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
sudo ovs-vsctl set open . external-ids:ovn-encap-ip=<loop back 以外の何らかのIP> # ループ計算でCPUを張り付かせる可能性があるので、loは回避したい。
sudo ovs-vsctl set open . external-ids:system-id=$(hostname)

sudo systemctl restart ovn-controller
```

## 注意: system-id設定前にovn-controllerが自動起動している場合がある
`apt install ovn-host`直後、ovn-controllerがsystem-id未設定のまま一度起動し、
ランダムUUIDでChassisをSB DBに登録してしまうことがある。
その後system-id等を設定してrestartしても、古いChassisレコードは自動削除されない。

対策: 一連の設定完了後、必ず以下で重複が無いか確認する。
```bash
sudo ovn-sbctl show
```
もし複数Chassisが見えたら、使われていない方を削除する。
```bash
sudo ovn-sbctl chassis-del <古いChassis名またはUUID>
```

Socketは、OVN界隈ではRootでしか触らないので、アプリから触る際にはEndpointが妥当。
```bash
sudo ovn-nbctl set-connection ptcp:6641:127.0.0.1
```
でTCP受付を行って、アプリエンドポイントもそこにする。

OVSも同様。
```bash
sudo ovs-vsctl set-manager ptcp:6640:127.0.0.1
```

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

Poolも、Libvirt入れただけだとできないので、
```bash
sudo virsh pool-define-as default dir \
  --target /var/lib/libvirt/images

sudo virsh pool-start default
sudo virsh pool-autostart default
```

あと、CloudeInitのために
```bash
sudo apt update
sudo apt install cloud-image-utils
```