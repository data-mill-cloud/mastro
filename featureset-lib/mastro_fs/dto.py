class FeatureSet:
    def __init__(self, name, version, description, labels = {}, features = []):
        self.name = name
        self.version = version
        self.description = description
        self.features = [Feature(**f) if isinstance(f, dict) else f for f in features] 
        self.labels = labels
        self.labels = labels

    def __eq__(self, other):
        return self.__dict__ == other.__dict__

class Feature:
    def __init__(self, name, value, data_type):
        self.name = name
        self.value = value
        self.data_type = data_type

    def __eq__(self, other):
        return self.__dict__ == other.__dict__