# This Vagrantfile defines 3 boxes:
# 1. build - download and build imagemagick and dependencies from source.
# 2. dev - a VM with bazel, to verify that builds work.
# 3. prod - a VM with nothing, to verify that running a binary works.
#
# During development, you will want to run rsync-auto to sync updates from
# host to guest.
#
# 	vagrant rsync-auto

MEMORY_MB = 4096
CPU_CORES = 1

Vagrant.configure("2") do |config|
  # Mitigate "VBoxHeadless + logd using all available CPU"
  $enable_serial_logging = false

  config.vagrant.plugins = ["vagrant-scp"]
  config.vm.network "private_network", type: "dhcp"
  config.vm.provider :virtualbox do |virtualbox, override|
    virtualbox.memory = MEMORY_MB
    virtualbox.cpus = CPU_CORES
    virtualbox.gui = false
  end

  # Builder box runs the download-build-imagemagick.sh, leaving the result in
  # /home/vagrant/imagemagick-build
  config.vm.define "build", autostart: false do |build|
    build.vm.box = "centos/6"
    build.vm.box_download_checksum_type = "sha256"
    build.vm.box_download_checksum = "68d96e36b94b5a2c57efc7790008c2c00b25c323f8e396160a98c5c7df6221ed"
    build.vm.box_url = "https://cloud.centos.org/centos/6/vagrant/x86_64/images/CentOS-6-x86_64-Vagrant-1812_01.VirtualBox.box"
    build.vm.provision "shell", inline: "yum -y install gcc gcc-c++"
    build.vm.provision "shell", privileged: false, path: "download-build-imagemagick.sh"
    build.vm.synced_folder ".", "/vagrant", disabled: true
  end

  # Dev box emulates a developer workstation for a library that uses this library.
  # They have bazel and the repo, but not ImageMagick.
  config.vm.define "dev" do |developer|
    developer.vm.box = "centos/7"
    developer.vm.box_download_checksum_type = "sha256"
    developer.vm.box_download_checksum = "b24c912b136d2aa9b7b94fc2689b2001c8d04280cf25983123e45b6a52693fb3"
    developer.vm.box_url = "https://cloud.centos.org/centos/7/vagrant/x86_64/images/CentOS-7-x86_64-Vagrant-1803_01.VirtualBox.box"
    developer.vm.provision "ansible" do |ansible|
      ansible.playbook = "vagrant-dev-playbook.yml"
      ansible.verbose = true
    end
    developer.vm.synced_folder ".", "/vagrant", type: "rsync",
      rsync__exclude: [
        "bazel-imagick/",
        "bazel-out/",
        "bazel-bin/",
        "bazel-genfiles/",
        "bazel-testlogs/",
      ]
  end

  # Prod box emulates a deployment target, with nothing installed.
  # Copying a binary here and running it serves to verify that it is sufficiently self-contained.
  config.vm.define "prod" do |prod|
    prod.vm.box = "centos/6"
    prod.vm.box_download_checksum_type = "sha256"
    prod.vm.box_download_checksum = "68d96e36b94b5a2c57efc7790008c2c00b25c323f8e396160a98c5c7df6221ed"
    prod.vm.box_url = "https://cloud.centos.org/centos/6/vagrant/x86_64/images/CentOS-6-x86_64-Vagrant-1812_01.VirtualBox.box"
    prod.vm.synced_folder ".", "/vagrant", disabled: true
    prod.vm.provider :virtualbox do |virtualbox, override|
      virtualbox.memory = 512
      virtualbox.cpus = 1
    end
  end

end
