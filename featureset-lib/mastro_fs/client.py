import http.client
import json
from mastro_fs.dto import FeatureSet, PaginatedFeatureSets

class MastroFeatureStoreClient:
    '''A client to interact with Mastro feature stores'''
    def __init__(self, host, port):
        self.host = host
        self.port = port
        self.fs_service = "{}:{}".format(host, port)

    def create_featureset(self, featureset):
        '''uploads the provided featureset on the remote store'''
        connection = http.client.HTTPConnection(self.host, self.port)
        headers = {'Content-type': 'application/json'}
        connection.request('PUT', '/featureset/', featureset.toJSON(), headers)
        response = connection.getresponse().read().decode()
        connection.close()
        return FeatureSet.fromJSON(response)

    def get_featureset_by_id(self, id):
        '''returns a featureset matching the provided id'''
        connection = http.client.HTTPConnection(self.host, self.port)
        connection.request("GET", "/featureset/id/{}".format(id))
        response = connection.getresponse().read().decode()
        connection.close()
        return FeatureSet.fromJSON(response)

    def get_featureset_by_name(self, name, limit=10, page=1):
        '''returns a featureset matching the provided name'''
        connection = http.client.HTTPConnection(self.host, self.port)
        connection.request("GET", "/featureset/name/{}?limit={}&page={}".format(name, limit, page))
        response = connection.getresponse().read().decode()
        connection.close()
        return PaginatedFeatureSets.fromJSON(response)

    def search_featuresets_by_labels(self, labels, limit=10, page=1):
        '''returns a paginated list of featureset that match the provided labels'''
        connection = http.client.HTTPConnection(self.host, self.port)
        headers = {'Content-type': 'application/json'}
        payload = {"labels": labels, "limit": limit, "page": page}
        connection.request('POST', '/labels', json.dumps(payload), headers)
        response = connection.getresponse().read().decode()
        connection.close()
        return PaginatedFeatureSets.fromJSON(response)

    def search(self, query, limit=10, page=1):
        connection = http.client.HTTPConnection(self.host, self.port)
        headers = {'Content-type': 'application/json'}
        payload = {"query": query, "limit": limit, "page": page}
        connection.request('POST', '/search', json.dumps(payload), headers)
        response = connection.getresponse().read().decode()
        connection.close()
        return PaginatedFeatureSets.fromJSON(response)

    def list_all(self, limit=10, page=1):
        connection = http.client.HTTPConnection(self.host, self.port)
        connection.request("GET", "/featureset/?limit={}&page={}".format(limit, page))
        response = connection.getresponse().read().decode()
        connection.close()
        return PaginatedFeatureSets.fromJSON(response)