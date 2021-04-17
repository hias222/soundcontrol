# install

## Preparation

RPI needed with

* npm installed
* nginx installed

## Ansible

ansible-playbook -i inventories/production/hosts testService.yml --limit=rpi 