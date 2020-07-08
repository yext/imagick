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
    build.vm.box = "ubuntu/bionic64"
    build.vm.box_version = "20200701.0.0"
    build.vm.provision "shell", inline: "apt update"
    build.vm.provision "shell", inline: "apt install -y bash build-essential"
    build.vm.provision "shell", privileged: false, path: "download-build-imagemagick.sh"
    build.vm.synced_folder ".", "/vagrant", disabled: true
  end

  # Dev box emulates a developer workstation for a library that uses this library.
  # They have bazel and the repo, but not ImageMagick.
  config.vm.define "dev" do |developer|
    developer.vm.box = "ubuntu/bionic64"
    developer.vm.box_version = "20200701.0.0"
    developer.vm.provision "shell", inline: "apt update"
    developer.vm.provision "shell", inline: "apt install -y python"
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
    prod.vm.box = "centos/7"
    prod.vm.box_version = "1809.1"
    prod.vm.synced_folder ".", "/vagrant", disabled: true
    prod.vm.provider :virtualbox do |virtualbox, override|
      virtualbox.memory = 512
      virtualbox.cpus = 1
    end
  end

end
