import json

class FeatureSet:
    def __init__(self, name, version, description = None, labels = {}, features = [], inserted_at = None):
        self.name = name
        self.version = version
        self.description = description
        self.features = [Feature(**f) if isinstance(f, dict) else f for f in features] 
        self.labels = labels
        self.labels = labels

    def __eq__(self, other):
        return self.__dict__ == other.__dict__

    def toJSON(self):
        return json.dumps(self, default=lambda o: o.__dict__, sort_keys=True)

    @staticmethod
    def fromJSON(json_string):
        j = json.loads(json_string)
        return FeatureSet(**j)

class Feature:
    def __init__(self, name, value, data_type):
        self.name = name
        self.value = value
        self.data_type = data_type

    def __eq__(self, other):
        return self.__dict__ == other.__dict__
    
    def toJSON(self):
        return json.dumps(self, default=lambda o: o.__dict__, sort_keys=True)

class PaginatedFeatureSets:
    def __init__(self, pagination, data):
        self.pagination = Pagination(**pagination) if isinstance(pagination, dict) else pagination
        self.data = [FeatureSet(**f) if isinstance(f, dict) else f for f in data]

    def __eq__(self, other):
        return self.__dict__ == other.__dict__

    def toJSON(self):
        return json.dumps(self, default=lambda o: o.__dict__, sort_keys=True)

    @staticmethod
    def fromJSON(json_string):
        j = json.loads(json_string)
        return PaginatedFeatureSets(**j)

class Pagination:
    def __init__(self, total, page, perPage, prev, next, totalPage):
        self.total = total
        self.page = page
        self.perPage = perPage
        self.prev = prev
        self.next = next
        self.totalPage = totalPage