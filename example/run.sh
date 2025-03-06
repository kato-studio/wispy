# 
# Fixed this by running "$ sysctl net.ipv4.ip_unprivileged_port_start=80"
# on local WSL env this should be properly configured on non-private envs
# -------
# sudo setcap 'cap_net_bind_service=ep cap_sys_admin=ep' ./tmp/main
# sleep 0.5
sudo sysctl net.ipv4.ip_unprivileged_port_start=80
sleep 0.5
air

