import {
  FieldOverrideContext,
  FieldType,
  getFieldDisplayName,
  PanelPlugin,
  ReducerID,
  standardEditorsRegistry,
  identityOverrideProcessor,
} from '@grafana/data';
import { TableFieldOptions, TableCellOptions, TableCellDisplayMode } from '@grafana/schema';

import { PaginationEditor } from './PaginationEditor';
import { TableCellOptionEditor } from './TableCellOptionEditor';
import { TablePanel } from './TablePanel';
import { tableMigrationHandler, tablePanelChangedHandler } from './migrations';
import { PanelOptions, defaultPanelOptions, defaultPanelFieldConfig } from './models.gen';
import { TableSuggestionsSupplier } from './suggestions';

const footerHeaderCategory = 'Table header & footer';
const cellCategory = ['Cell Options'];

export const plugin = new PanelPlugin<PanelOptions, TableFieldOptions>(TablePanel)
  .setPanelChangeHandler(tablePanelChangedHandler)
  .setMigrationHandler(tableMigrationHandler)
  .useFieldConfig({
    useCustomConfig: (builder) => {
      builder
        .addNumberInput({
          path: 'minWidth',
          name: 'Minimum column width',
          description: 'The minimum width for column auto resizing',
          settings: {
            placeholder: '150',
            min: 50,
            max: 500,
          },
          shouldApply: () => true,
          defaultValue: defaultPanelFieldConfig.minWidth,
        })
        .addNumberInput({
          path: 'width',
          name: 'Column width',
          settings: {
            placeholder: 'auto',
            min: 20,
            max: 300,
          },
          shouldApply: () => true,
          defaultValue: defaultPanelFieldConfig.width,
        })
        .addRadio({
          path: 'align',
          name: 'Column alignment',
          settings: {
            options: [
              { label: 'auto', value: 'auto' },
              { label: 'left', value: 'left' },
              { label: 'center', value: 'center' },
              { label: 'right', value: 'right' },
            ],
          },
          defaultValue: defaultPanelFieldConfig.align,
        })
        .addCustomEditor<void, TableCellOptions>({
          id: 'cellOptions',
          path: 'cellOptions',
          name: 'Cell Type',
          editor: TableCellOptionEditor,
          override: TableCellOptionEditor,
          defaultValue: defaultPanelFieldConfig.cellOptions,
          process: identityOverrideProcessor,
          category: cellCategory,
          shouldApply: () => true,
        })
        .addBooleanSwitch({
          path: 'inspect',
          name: 'Cell value inspect',
          description: 'Enable cell value inspection in a modal window',
          defaultValue: false,
          category: cellCategory,
          showIf: (cfg) => {
            return (
              cfg.cellOptions.type === TableCellDisplayMode.Auto ||
              cfg.cellOptions.type === TableCellDisplayMode.JSONView ||
              cfg.cellOptions.type === TableCellDisplayMode.ColorText ||
              cfg.cellOptions.type === TableCellDisplayMode.ColorBackground
            );
          },
        })
        .addBooleanSwitch({
          path: 'filterable',
          name: 'Column filter',
          description: 'Enables/disables field filters in table',
          defaultValue: defaultPanelFieldConfig.filterable,
        })
        .addBooleanSwitch({
          path: 'hidden',
          name: 'Hide in table',
          defaultValue: undefined,
          hideFromDefaults: true,
        });
    },
  })
  .setPanelOptions((builder) => {
    builder
      .addBooleanSwitch({
        path: 'showHeader',
        category: [footerHeaderCategory],
        name: 'Show table header',
        defaultValue: defaultPanelOptions.showHeader,
      })
      .addBooleanSwitch({
        path: 'showRowNums',
        name: 'Show row numbers',
        defaultValue: defaultPanelOptions.showRowNums,
      })
      .addBooleanSwitch({
        path: 'footer.show',
        category: [footerHeaderCategory],
        name: 'Show table footer',
        defaultValue: defaultPanelOptions.footer?.show,
      })
      .addCustomEditor({
        id: 'footer.reducer',
        category: [footerHeaderCategory],
        path: 'footer.reducer',
        name: 'Calculation',
        description: 'Choose a reducer function / calculation',
        editor: standardEditorsRegistry.get('stats-picker').editor as any,
        defaultValue: [ReducerID.sum],
        showIf: (cfg) => cfg.footer?.show,
      })
      .addBooleanSwitch({
        path: 'footer.countRows',
        category: [footerHeaderCategory],
        name: 'Count rows',
        description: 'Display a single count for all data rows',
        defaultValue: defaultPanelOptions.footer?.countRows,
        showIf: (cfg) => cfg.footer?.reducer?.length === 1 && cfg.footer?.reducer[0] === ReducerID.count,
      })
      .addMultiSelect({
        path: 'footer.fields',
        category: [footerHeaderCategory],
        name: 'Fields',
        description: 'Select the fields that should be calculated',
        settings: {
          allowCustomValue: false,
          options: [],
          placeholder: 'All Numeric Fields',
          getOptions: async (context: FieldOverrideContext) => {
            const options = [];
            if (context && context.data && context.data.length > 0) {
              const frame = context.data[0];
              for (const field of frame.fields) {
                if (field.type === FieldType.number) {
                  const name = getFieldDisplayName(field, frame, context.data);
                  const value = field.name;
                  options.push({ value, label: name } as any);
                }
              }
            }
            return options;
          },
        },
        defaultValue: '',
        showIf: (cfg) =>
          (cfg.footer?.show && !cfg.footer?.countRows) ||
          (cfg.footer?.reducer?.length === 1 && cfg.footer?.reducer[0] !== ReducerID.count),
      })
      .addCustomEditor({
        id: 'footer.enablePagination',
        path: 'footer.enablePagination',
        name: 'Enable pagination',
        editor: PaginationEditor,
      });
  })
  .setSuggestionsSupplier(new TableSuggestionsSupplier());
