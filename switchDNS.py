import requests
import time
import dns.resolver
from prettytable import PrettyTable
import func_timeout
from ping3 import ping


dns_list = [
    "114.114.114.114",
    "8.8.8.8",
    "8.8.4.4",
    "218.102.23.228",
    "211.136.192.6",
    "223.5.5.5",
    "168.126.63.1",
    "168.126.63.2",
    "168.95.1.1",
    "168.95.192.1",
    "203.80.96.9",
    "61.10.0.130",
    "61.10.1.130",
    "208.67.222.222",
    "208.67.220.220",
    "202.14.67.4",
    "203.80.96.10",
    "202.14.67.14",
    "198.153.194.1",
    "198.153.192.1",
    "112.106.53.22",
    "168.126.63.1",
    "168.95.192.1",
    "198.153.194.1",
    "210.2.4.8",
    "203.80.96.9",
    "220.67.240.221",
    "84.200.69.80",
    "81.218.119.11",
    "180.76.76.76",
    "119.29.29.29",
]


def single_dns_test(name_server):
    '''
        return reachable, ping, latency, dl_speed, up_speed
    '''
    reachable, latency, dl_addr, up_addr = qname_avalible(name_server)
    if reachable == 1:
        ping = ping_latency(name_server)
        dl_speed, up_speed = speedtest(dl_addr, up_addr)
        return name_server, reachable, ping, latency, dl_speed, up_speed
    elif reachable == 0:
        return name_server, reachable, 0, 0, 0, 0
    elif reachable == -1:
        return name_server, reachable, 0, 0, 0, 0


def qname_avalible(name_server):
    '''
        return reachable, resolve_latency, dl_addr, up_addr

        reachable == 0  # cannot resolve
                == -1   # notreachable
                == 1    # resolve success
    '''
    Nintendo_Dl_Url = "ctest-dl-lp1.cdn.nintendo.net"
    Nintendo_Up_Url = "ctest-ul-lp1.cdn.nintendo.net"
    try:
        Resolver = dns.resolver.Resolver()
        Resolver.nameservers = [name_server]
        t1 = time.time()
        Dl_result = Resolver.resolve(Nintendo_Dl_Url, "A")
        Up_result = Resolver.resolve(Nintendo_Up_Url, "A")
        t2 = time.time()
    except dns.resolver.NXDOMAIN:
        return 0, 0, '', ''
    except:
        return -1, 0, '', ''
    else:
        return 1, int((t2-t1)*1000), Dl_result.rrset[0].address, Up_result.rrset[0].address


def ping_latency(ip_addr):
    '''
        return average ping latency
    '''
    _ = ping(ip_addr)
    return int(_*1000) if _ else -1


def speedtest(dl_addr, up_addr, dl_timeout=15, up_timeout=5):
    '''
        return download_speed, upload_speed # MB/s
    '''
    headers_ua = {
        "user-agent": "Nintendo NX",
        "host": "ctest-dl-lp1.cdn.nintendo.net"
    }
    @func_timeout.func_set_timeout(dl_timeout)
    def dl_req(dl_addr, headers_ua):
        _ = requests.get("http://" + dl_addr + "/30m",
                    headers=headers_ua)
    @func_timeout.func_set_timeout(up_timeout)
    def up_req(up_addr, headers_ua):
        _ = requests.post("http://" + up_addr + "/1m",
                            data={" ": " " * 1024**2}, headers=headers_ua)

    try:
        T1 = time.time()
        dl_req(dl_addr, headers_ua)
    except func_timeout.exceptions.FunctionTimedOut:
        dl_speed = 'too_slow'
    else:
        dl_speed = round(30 / (time.time() - T1), 2)

    try:
        T3 = time.time()
        up_req(up_addr, headers_ua)
    except func_timeout.exceptions.FunctionTimedOut:
        up_speed = 'too_slow'
    else:
        up_speed = round(1 / (time.time() - T3), 2)

    return dl_speed, up_speed


def thread_control(dns_list):
    # with concurrent.futures.ThreadPoolExecutor() as executor:
    #     to_do = []
    #     for name_server in dns_list:
    #         future = executor.submit(single_dns_test, name_server)
    #         to_do.append(future)
    # return [_.result() for _ in to_do]
    _ = []
    for name_server in dns_list:
        _.append(single_dns_test(name_server))
    return _


def table_print(test_rusults):
    x = PrettyTable()
    x.field_names = ["nameserver", "reachable", "ping",
                     "latency", "dl_speed", "up_speed"]
    [ x.add_row(i) for i in test_rusults if (i[1] == 1) ]
    print(x)

_ = time.time()
table_print(thread_control(dns_list))
print(time.time()- _)
input()