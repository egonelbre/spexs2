===============
SPeXS Changelog
===============

spxs v0.9
=========

changes
-------

* in configuration alphabet.groupName.group is now alphabet.groupName.elements
* removed characters definition from configuration, it will be automatically inferred from the input
* added alphabet.separator for alternative inputs. The lines and group elements will be broken based on that separator. See data/text.json for an alternative example.

fixes
-----

* fixed append problem introduced in v0.7@82a18155
