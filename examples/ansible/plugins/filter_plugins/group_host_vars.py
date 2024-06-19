#!/usr/bin/python

from ansible import errors

class FilterModule(object):
    def filters(self):
        return {
            'each_hostvars': self.each_hostvars,
            'keep_attribute': self.keep_attribute,
            'each_format': self.each_format
        }

    def each_hostvars(self, hostnames, hostvars):
        ''' Get a list of hostvars from a list of hostnames '''
        try:
            return [
                hostvars[hostname]
                for hostname in hostnames
            ]
        except Exception as e:
            raise errors.AnsibleFilterError(
                    'each_hostvars plugin error: {0}, hostnames={1}'.format(str(e), str(hostnames))
            )

    def keep_attribute(self, server_hostvars, key):
        ''' Get a list of values corresponding to a key from a list of hostvars '''
        try:
            return [
                hostvars[key]
                for hostvars in server_hostvars
            ]
        except Exception as e:
            raise errors.AnsibleFilterError(
                    'keep_attribute plugin error: {0}, key={1}'.format(str(e), str(key))
            )

    def each_format(self, server_hostvars, string_format):
        ''' Get a list of formated strings using a list of hostvars '''
        try:
            return [
                string_format.format_map(hostvars)
                for hostvars in server_hostvars
            ]
        except Exception as e:
            raise errors.AnsibleFilterError(
                    'each_format plugin error: {0}, string_format={1}'.format(str(e), str(string_format))
            )
