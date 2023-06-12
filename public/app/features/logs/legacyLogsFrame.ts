import { DataFrame, FieldCache, FieldType, Field } from '@grafana/data';

import type { LogsFrame, Attributes } from './logsFrame';

function getAttributes(frame: DataFrame, cache: FieldCache, lineField: Field): Attributes[] | undefined {
  const useLabelsField = frame.meta?.custom?.frameType === 'LabeledTimeValues';

  if (!useLabelsField) {
    const lineLabels = lineField.labels;
    if (lineLabels != null) {
      const result = new Array(frame.length);
      result.fill(lineLabels);
      return result;
    } else {
      return undefined;
    }
  }

  const labelsField = cache.getFieldByName('labels');

  if (labelsField == null) {
    return undefined;
  }

  return labelsField.values;
}

export function parseLegacyLogsFrame(frame: DataFrame): LogsFrame | null {
  const cache = new FieldCache(frame);
  const timeField = cache.getFields(FieldType.time)[0];
  const bodyField = cache.getFields(FieldType.string)[0];

  // these two are mandatory
  if (timeField === undefined || bodyField === undefined) {
    return null;
  }

  const timeNanosecondField = cache.getFieldByName('tsNs');
  const severityField = cache.getFieldByName('level');
  const idField = cache.getFieldByName('id');

  return {
    timeField,
    bodyField,
    timeNanosecondField,
    severityField,
    idField,
    attributes: getAttributes(frame, cache, bodyField),
  };
}
