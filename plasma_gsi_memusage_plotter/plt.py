import sys
from matplotlib.ticker import FuncFormatter
import matplotlib.pyplot as plt
import numpy as np

import matplotlib.cbook as cbook

filename = sys.argv[1]

data = cbook.get_sample_data(filename, asfileobj=False)
plt.plotfile(data, ('time', 'rss', 'storage_used', 'je_allocated', 'je_resident', 'reclaim_pending', 'queue_size'), subplots=False, delimiter=',')

plt.xlabel(r'seconds')
plt.ylabel(r'memory_used')
yaxis = plt.twiny()
yaxis.yaxis.set_major_formatter(FuncFormatter(lambda y,_: '%1.f' % (y/(1024*1024))))

plt.plotfile(data, ('time', 'mainrecord', 'docrecords'), subplots=False, delimiter=',')
plt.show()


