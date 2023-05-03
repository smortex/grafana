import { MetricKind } from '../types/query';
import { MetricDescriptor, ValueTypes } from '../types/types';

export const createMockMetricDescriptor = (overrides?: Partial<MetricDescriptor>): MetricDescriptor => {
  return {
    metricKind: MetricKind.CUMULATIVE,
    valueType: ValueTypes.DOUBLE,
    type: 'type',
    unit: 'unit',
    service: 'service',
    serviceShortName: 'srv',
    displayName: 'metricName',
    description: 'description',
    ...overrides,
  };
};
