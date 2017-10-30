from locust import HttpLocust
from locust import TaskSet
from locust import task


class PingTaskSet(TaskSet):
    min_wait = 10
    max_wait = 10

    @task
    def my_task(self):
        host_header = 'go.counter.examples.local.appcelerator.io'
        ssl_verify = False
        self.client.get('/ping', verify=ssl_verify, headers={'Host': host_header})


class PingLocust(HttpLocust):
    task_set = PingTaskSet
