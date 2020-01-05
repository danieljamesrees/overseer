Experimental
============

A block chain used to store experiment results and analyse them using AI.

Requirements
------------

If I was actually going to do something useful, I'd probably set up the block chain on AWS (not having any access to non-local computing resources myself). As such, the role would likely require the boto package.

Role Variables
--------------

name: The name of the block chain - MANDATORY
number_of_replicas: The number of Deployment replicates to create

is_public: Whether the block chain is publicly accessible or not 
block_time: The time (in seconds) to generate one extra block in the chain

Dependencies
------------

No roles from Galaxy are required.

No parameters need to be set for any other roles.

No variables are used from other roles.

Example Playbook
----------------

Use this role as in the following example:

    - hosts: servers
      roles:
         - { role: experimental, name: mad-science, block_time: 60 }

License
-------

BSD

Author Information
------------------

Daniel Rees
