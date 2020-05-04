import json
import logging
import os
from enum import Enum
from io import StringIO

import jmespath
import yaml
from jinja2 import (Environment, FileSystemLoader, TemplateNotFound,
                    select_autoescape)

from graviteeio_cli.exeptions import GraviteeioError
from graviteeio_cli.graviteeio.extensions.jinja_filters import filter_loader
from graviteeio_cli.graviteeio.apim.apis.utils import update_dic_with_set

from ..... import environments
from . import init_template as init

logger = logging.getLogger("class-ApiSchema")


def load_yaml(stream):
    return yaml.load(stream, Loader=yaml.SafeLoader)

def load_json(stream):
    io = StringIO(stream)
    return json.load(io)

def dump_yaml(data):
    return yaml.dump(data)

def dump_json(data):
    io = StringIO()
    json.dump(data, io, indent=2)
    return io.getvalue()

map_extention = {}
class Data_Template_Format(Enum):
    YAML = {
        'num': 1,
        'extentions': ['.yml', '.yaml'],
        'load': load_yaml,
        'dump': dump_yaml
    }
    JSON = {
        'num': 2,
        'extentions': ['.json'],
        'load': load_json,
        'dump': dump_json
    }

    def __init__(self, values):
        self.num = values['num']
        self.extentions = values['extentions']
        self.load = values['load']
        self.dump = values['dump']

        for extention in values['extentions']:
            map_extention[extention] = self
    
    @staticmethod
    def find(extention):
        to_return = None
        if extention in map_extention:
            to_return = map_extention[extention]
        return to_return

    @staticmethod
    def list_name():
        return list (map (lambda c: c.name, Data_Template_Format))
    
    @staticmethod
    def value_of(value):
        for data_type in Data_Template_Format:
            if data_type.name == value.upper():
                return data_type

    @staticmethod
    def extention_list():
        tuple_to_return = []

        for format in Data_Template_Format:
            if type(format.extentions) is tuple:
                for extention in format.extentions:
                    tuple_to_return.append(extention)
            else:
                tuple_to_return.append(format.extentions)

        return tuple(tuple_to_return)


class ApiSchema:

    def __init__(self, resources_folder, value_file = None):
        self.folders = {}
        self.files = {}
        self.loaded_schema = False

        self.folders["templates_folder"] = "{}/{}".format(resources_folder, environments.APIM_API_TEMPLATES_FOLDER)
        self.folders["settings_folder"] = "{}/{}".format(resources_folder, environments.APIM_API_SETTING_FOLDER)

        self.files["root_template_path_file"] = "{}/{}".format(self.folders["templates_folder"], environments.APIM_API_TEMPLATE_FILE)

        if not value_file or len(value_file) == 0:
            value_file = '{}/{}'.format(resources_folder, environments.APIM_API_VALUE_FILE_NAME)

        self.files["value_file"] = value_file


    def generate_schema(self, format = Data_Template_Format.YAML, api_def = None, debug = False):
        for key in self.folders:
            if debug:
                print("mkdir {}".format(self.folders[key]))
            else:
                try:
                    os.mkdir(self.folders[key])
                except OSError:
                    print ("Creation of the file %s failed" % self.folders[key])
                else:
                    print ("Successfully created directory %s " % self.folders[key])

        write_files = {
            self.files["root_template_path_file"].format(format.extentions[0]): init.templates[format.name.lower()]["template"],
            self.files["value_file"].format(format.extentions[0]): init.templates[format.name.lower()]["value_file"],
            "{}/{}".format(self.folders["settings_folder"],"Http{}".format(format.extentions[0])): init.templates[format.name.lower()]["setting_http"]
        }

        if api_def:
            write_files = self._api_def_to_template(api_def, format)
       
        for key in write_files:
            if debug:
                print("write file {}".format(key))
            else:
                try:
                    with open(key, 'x') as f:
                        f.write(write_files[key])
                except OSError:
                    print ("Creation of the file %s failed" % key)
                else:
                    print ("Successfully created file %s " % key)

    def _api_def_to_template(self, api_def , format):
        values = {}

        values["version"]= api_def["version"]
        values["name"]= api_def["name"]
        values["description"] = api_def["description"]

        api_def["version"] = "{{ Values.version}}"
        api_def["name"] = "{{ Values.name}}"
        api_def["description"] = "{{ Values.description}}"

        write_files = {
            self.files["root_template_path_file"].format(format.extentions[0]): format.dump(api_def),
            self.files["value_file"].format(format.extentions[0]): format.dump(values),
            "{}/{}".format(self.folders["settings_folder"],"Http{}".format(format.extentions[0])): ""
        }
        
        return write_files

    def get_api_data(self, debug = None, set_values = []):
        if not self.loaded_schema: self._load_schema()

        api_data_rendered = None

        for set_value in set_values:
            self.api_vars["Values"] = update_dic_with_set(set_value, self.api_vars["Values"])

        api_data_rendered = self.template.render(self.api_vars)

        if debug:
            print("Render:")
            print(api_data_rendered)
        
        api_data_dic =  self.template_format.load(api_data_rendered)

        if 'version' in api_data_dic:
            api_data_dic['version'] = str(api_data_dic['version'])

        return api_data_dic

    def _load_schema(self):
        for key in self.folders:
            if not os.path.exists(self.folders[key]):
                raise GraviteeioError("Missing folder {}".format(self.folders[key]))

        root_template_file = None
        
        value_file = None
        if os.path.exists(self.files["value_file"]):
            value_file = self.files["value_file"]

        root_template_path_file = self.files["root_template_path_file"]
        for (data_format, extention) in ((data_format, extention) for data_format in Data_Template_Format for extention in data_format.extentions):
            if not root_template_file and os.path.exists(root_template_path_file.format(extention)):
                self.template_format = data_format;
                root_template_file = environments.APIM_API_TEMPLATE_FILE.format(extention)
            
            file = self.files["value_file"].format(extention)
            if not value_file and os.path.exists(file):
                value_file_format = data_format
                value_file = file
            if root_template_file and value_file:
                break;

        templates_folder = self.folders["templates_folder"]
        j2_env = Environment(loader = FileSystemLoader(templates_folder), trim_blocks=False, autoescape=False)
        filter_loader(j2_env)

        try:
            template = j2_env.get_template(root_template_file)
        except TemplateNotFound:
            raise GraviteeioError("Template not found, try to load {}".format(root_template_file))
        self.template = template

        try:
            with open(value_file, 'r') as f:
                api_value_string = f.read()
        except FileNotFoundError:
            raise GraviteeioError("Missing values file {}".format(value_file))

        self.api_vars = {}
        self.api_vars["Values"] = value_file_format.load(api_value_string)

        settings_folder = self.folders["settings_folder"]
        config_files = []

        for file in os.listdir(settings_folder):
            if not file.startswith(('_', ".")):
                try:
                    with open("/".join([settings_folder, file]), 'r') as f:
                        config_string = f.read()
                except FileNotFoundError:
                    raise GraviteeioError("No such file {}".format(file))
                
                filename, file_extension = os.path.splitext(file)
                file_format = Data_Template_Format.find(file_extension)

                if file_format:
                    self.api_vars[filename] = file_format.load(config_string)
        
        self.loaded_schema = True
