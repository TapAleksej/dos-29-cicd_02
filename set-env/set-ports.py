import socket
from ruamel.yaml import YAML

yaml = YAML()
int_port = "8080"

def find_free_port(start, end):
    for port in range(start, end + 1):
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            if s.connect_ex(('localhost', port)) != 0:
                return port
    return None

def update_compose_ports(file_path, service_name, new_port):
    with open(file_path, 'r', encoding='utf-8') as f:
        data = yaml.load(f)

    service = data.get('services', {}).get(service_name)
    if service and 'ports' in service:
        ports = service['ports']
        if isinstance(ports, list) and ports:
            ports[0] = f"{new_port}:{new_port}"
        else:
            service['ports'] = [f"{new_port}:{int_port}"]
    else:
        if 'services' in data and service_name in data['services']:
            data['services'][service_name]['ports'] = [f"{new_port}:{int_port}"]
        else:
            raise ValueError(f"Сервис {service_name} не найден в файле")

    with open(file_path, 'w', encoding='utf-8') as f:
        yaml.dump(data, f)

if __name__ == "__main__":
    start_port = 8080
    end_port = 8090
    compose_file = '../compose.yml'
    service = 'api'

    free_port = find_free_port(start_port, end_port)
    if free_port is None:
        print("Свободных портов в заданном диапазоне не найдено")
    else:
        print(f"Найден свободный порт: {free_port}")
        update_compose_ports(compose_file, service, free_port)
        print(f"Файл {compose_file} обновлен. Порт для сервиса '{service}' установлен в {free_port}")

