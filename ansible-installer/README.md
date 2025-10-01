# Ansible Installer for CloudCurio Monorepo

This project provides Ansible playbooks and roles to install and configure various tools, databases, and Apache software suite on Linux baremetal.

## Structure
- `inventories/`: Host configurations for local and production.
- `group_vars/`: Variables to toggle installations.
- `roles/`: Individual roles for each software.
- `playbooks/`: Main playbooks.
- `tui/`: Go-based TUI using Bubble Tea to select and run installations.

## Usage
1. Edit `group_vars/all.yml` to toggle installs.
2. For local: `ansible-playbook -i inventories/local/hosts.yml playbooks/install-all.yml`
3. For remote: Update production/hosts.yml with SSH details, run similarly.
4. TUI: `cd tui && go run main.go` - Select options, then it will run Ansible.

## Notes
- Assumes Ubuntu/Debian.
- Some tools (e.g., Sentry, Langfuse) are Docker-oriented; adapted where possible.
- Apache suite includes major projects; extend roles for more (Ant, Maven, etc.).
- For full Apache list, add roles like apache-maven, apache-ant, etc., following similar patterns.
- TUI is basic; extend to update vars and exec ansible-playbook with --extra-vars.

To add more Apache projects:
- Create role e.g., `roles/apache-maven/tasks/main.yml` with download/extract/setup.
