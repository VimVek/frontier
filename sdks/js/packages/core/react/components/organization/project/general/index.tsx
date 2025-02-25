import {
  Button,
  Flex,
  InputField,
  Separator,
  Text,
  TextField
} from '@raystack/apsara';

import { yupResolver } from '@hookform/resolvers/yup';
import { useNavigate, useParams } from '@tanstack/react-router';
import { useEffect, useMemo } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { toast } from 'sonner';
import * as yup from 'yup';
import { useFrontier } from '~/react/contexts/FrontierContext';
import { usePermissions } from '~/react/hooks/usePermissions';
import { V1Beta1Organization, V1Beta1Project } from '~/src';
import { PERMISSIONS, shouldShowComponent } from '~/utils';

const projectSchema = yup
  .object({
    title: yup.string().required(),
    name: yup.string().required()
  })
  .required();

interface GeneralProjectProps {
  project?: V1Beta1Project;
  organization?: V1Beta1Organization;
}

export const General = ({ organization, project }: GeneralProjectProps) => {
  const {
    reset,
    control,
    handleSubmit,
    formState: { errors, isSubmitting }
  } = useForm({
    resolver: yupResolver(projectSchema)
  });
  let { projectId } = useParams({ from: '/projects/$projectId' });
  const { client } = useFrontier();

  useEffect(() => {
    reset(project);
  }, [reset, project]);

  const resource = `app/project:${projectId}`;
  const listOfPermissionsToCheck = [
    {
      permission: PERMISSIONS.UpdatePermission,
      resource
    },
    {
      permission: PERMISSIONS.DeletePermission,
      resource
    }
  ];

  const { permissions } = usePermissions(listOfPermissionsToCheck, !!projectId);

  const { canUpdateProject, canDeleteProject } = useMemo(() => {
    return {
      canUpdateProject: shouldShowComponent(
        permissions,
        `${PERMISSIONS.UpdatePermission}::${resource}`
      ),
      canDeleteProject: shouldShowComponent(
        permissions,
        `${PERMISSIONS.DeletePermission}::${resource}`
      )
    };
  }, [permissions, resource]);

  async function onSubmit(data: any) {
    if (!client) return;
    if (!organization?.id) return;
    if (!projectId) return;

    try {
      await client.frontierServiceUpdateProject(projectId, data);
      toast.success('Project updated');
    } catch ({ error }: any) {
      toast.error('Something went wrong', {
        description: error.message
      });
    }
  }

  return (
    <Flex direction="column" gap="large" style={{ paddingTop: '32px' }}>
      <form onSubmit={handleSubmit(onSubmit)}>
        <Flex direction="column" gap="medium" style={{ maxWidth: '320px' }}>
          <InputField label="Project title">
            <Controller
              render={({ field }) => (
                <TextField
                  {...field}
                  // @ts-ignore
                  size="medium"
                  placeholder="Provide project title"
                />
              )}
              control={control}
              name="title"
            />

            <Text size={1} style={{ color: 'var(--foreground-danger)' }}>
              {errors.title && String(errors.title?.message)}
            </Text>
          </InputField>
          <InputField label="Project name">
            <Controller
              render={({ field }) => (
                <TextField
                  {...field}
                  // @ts-ignore
                  size="medium"
                  disabled
                  placeholder="Provide project name"
                />
              )}
              control={control}
              name="name"
            />

            <Text size={1} style={{ color: 'var(--foreground-danger)' }}>
              {errors.title && String(errors.title?.message)}
            </Text>
          </InputField>
          {canUpdateProject ? (
            <Button variant="primary" size="medium" type="submit">
              {isSubmitting ? 'updating...' : 'Update project'}
            </Button>
          ) : null}
        </Flex>
      </form>
      <Separator />
      {canDeleteProject ? (
        <>
          <GeneralDeleteProject organization={organization} />
          <Separator />
        </>
      ) : null}
    </Flex>
  );
};

export const GeneralDeleteProject = ({ organization }: GeneralProjectProps) => {
  const { client } = useFrontier();
  let { projectId } = useParams({ from: '/projects/$projectId' });
  const navigate = useNavigate({ from: '/projects/$projectId' });
  const {
    handleSubmit,
    formState: { errors, isSubmitting }
  } = useForm();

  return (
    <Flex direction="column" gap="medium">
      <Text size={3} style={{ color: 'var(--foreground-muted)' }}>
        If you want to permanently delete this project and all of its data.
      </Text>
      <Button
        variant="danger"
        type="submit"
        size="medium"
        onClick={() =>
          navigate({
            to: `/projects/$projectId/delete`,
            params: { projectId: projectId }
          })
        }
      >
        {isSubmitting ? 'deleting...' : 'Delete project'}
      </Button>
    </Flex>
  );
};
