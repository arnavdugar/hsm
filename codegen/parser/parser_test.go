package parser_test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/arnavdugar/hsm/codegen/config"
	"github.com/arnavdugar/hsm/codegen/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed test/empty_action_name.yaml
var EmptyActionNameConfig string

func TestParseEmptyActionName(t *testing.T) {
	reader := strings.NewReader(EmptyActionNameConfig)

	_, err := parser.Parse(reader)

	require.EqualError(t, err, "action name must not be empty")
}

//go:embed test/duplicate_action.yaml
var DuplicateActionConfig string

func TestParseDuplicateAction(t *testing.T) {
	reader := strings.NewReader(DuplicateActionConfig)

	_, err := parser.Parse(reader)

	require.Error(t, err)
	assert.Equal(t, "duplicate action name: ActionName", err.Error())
}

//go:embed test/empty_state_name.yaml
var EmptyStateNameConfig string

func TestParseEmptyStateName(t *testing.T) {
	reader := strings.NewReader(EmptyStateNameConfig)

	_, err := parser.Parse(reader)

	require.EqualError(t, err, "state name must not be empty")
}

//go:embed test/duplicate_state.yaml
var DuplicateStateConfig string

func TestParseDuplicateState(t *testing.T) {
	reader := strings.NewReader(DuplicateStateConfig)

	_, err := parser.Parse(reader)

	require.Error(t, err)
	assert.Equal(t, "duplicate state name: StateName", err.Error())
}

//go:embed test/unknown_destination.yaml
var UnknownDestinationConfig string

func TestParseUnknownDestination(t *testing.T) {
	reader := strings.NewReader(UnknownDestinationConfig)

	_, err := parser.Parse(reader)

	require.Error(t, err)
	assert.Equal(t, `unknown destination "UnknownState" in "ActionName" transition from state "StateName"`, err.Error())
}

//go:embed test/unknown_transition_action.yaml
var UnknownTransitionActionConfig string

func TestParseUnknownTransitionAction(t *testing.T) {
	reader := strings.NewReader(UnknownTransitionActionConfig)

	_, err := parser.Parse(reader)

	require.Error(t, err)
	assert.Equal(t, `unknown action "UnknownActionName" in transitions for state "StateName"`, err.Error())
}

//go:embed test/duplicate_group.yaml
var DuplicateGroupConfig string

func TestParseDuplicateGroup(t *testing.T) {
	reader := strings.NewReader(DuplicateGroupConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, "duplicate group name: G")
}

//go:embed test/shared_namespace.yaml
var SharedNamespaceConfig string

func TestParseSharedNamespace(t *testing.T) {
	reader := strings.NewReader(SharedNamespaceConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, "name used by both a state and a group: A")
}

//go:embed test/empty_group_name.yaml
var EmptyGroupNameConfig string

func TestParseEmptyGroupName(t *testing.T) {
	reader := strings.NewReader(EmptyGroupNameConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, "group name must not be empty")
}

//go:embed test/unknown_subgroup.yaml
var UnknownSubgroupConfig string

func TestParseUnknownSubgroup(t *testing.T) {
	reader := strings.NewReader(UnknownSubgroupConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `unknown group "Missing" in group "G"`)
}

//go:embed test/unknown_member_state.yaml
var UnknownMemberStateConfig string

func TestParseUnknownMemberState(t *testing.T) {
	reader := strings.NewReader(UnknownMemberStateConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `unknown state "Missing" in group "G"`)
}

//go:embed test/state_used_as_subgroup.yaml
var StateUsedAsSubgroupConfig string

func TestParseStateUsedAsSubgroup(t *testing.T) {
	reader := strings.NewReader(StateUsedAsSubgroupConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `unknown group "A" in group "G"`)
}

//go:embed test/self_containment.yaml
var SelfContainmentConfig string

func TestParseSelfContainment(t *testing.T) {
	reader := strings.NewReader(SelfContainmentConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, "group containment cycle: G -> G")
}

//go:embed test/containment_cycle.yaml
var ContainmentCycleConfig string

func TestParseContainmentCycle(t *testing.T) {
	reader := strings.NewReader(ContainmentCycleConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, "group containment cycle: G -> H -> G")
}

//go:embed test/subgroup_before_parent.yaml
var SubgroupBeforeParentConfig string

func TestParseSubgroupBeforeParent(t *testing.T) {
	reader := strings.NewReader(SubgroupBeforeParentConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `group "G" must be declared before group "H"`)
}

//go:embed test/unknown_initial.yaml
var UnknownInitialConfig string

func TestParseUnknownInitial(t *testing.T) {
	reader := strings.NewReader(UnknownInitialConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `unknown initial target "Missing" in group "G"`)
}

//go:embed test/initial_state_outside_group.yaml
var InitialStateOutsideGroupConfig string

func TestParseInitialStateOutsideGroup(t *testing.T) {
	reader := strings.NewReader(InitialStateOutsideGroupConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `initial state "B" is not a member of group "G"`)
}

//go:embed test/initial_overlapping_group_is_not_a_subgroup.yaml
var InitialOverlappingGroupIsNotASubgroupConfig string

func TestParseInitialOverlappingGroupIsNotASubgroup(t *testing.T) {
	reader := strings.NewReader(InitialOverlappingGroupIsNotASubgroupConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `initial group "H" is not a subgroup of group "G"`)
}

//go:embed test/initial_points_to_self.yaml
var InitialPointsToSelfConfig string

func TestParseInitialPointsToSelf(t *testing.T) {
	reader := strings.NewReader(InitialPointsToSelfConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `initial group "G" is not a subgroup of group "G"`)
}

//go:embed test/initial_cycle.yaml
var InitialCycleConfig string

func TestParseInitialCycle(t *testing.T) {
	reader := strings.NewReader(InitialCycleConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `initial group "G" is not a subgroup of group "H"`)
}

//go:embed test/initial_subgroup_missing_initial.yaml
var InitialSubgroupMissingInitialConfig string

func TestParseInitialSubgroupMissingInitial(t *testing.T) {
	reader := strings.NewReader(InitialSubgroupMissingInitialConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `initial value for group "G": group "H" has no initial value`)
}

//go:embed test/destination_group_missing_initial.yaml
var DestinationGroupMissingInitialConfig string

func TestParseDestinationGroupMissingInitial(t *testing.T) {
	reader := strings.NewReader(DestinationGroupMissingInitialConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `group "G" has no initial value in "Go" transition from state "A"`)
}

//go:embed test/unknown_group_action.yaml
var UnknownGroupActionConfig string

func TestParseUnknownGroupAction(t *testing.T) {
	reader := strings.NewReader(UnknownGroupActionConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `unknown action "Missing" in transitions for group "G"`)
}

//go:embed test/unknown_group_action_destination.yaml
var UnknownGroupActionDestinationConfig string

func TestParseUnknownGroupActionDestination(t *testing.T) {
	reader := strings.NewReader(UnknownGroupActionDestinationConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `unknown destination "Missing" in "Go" transition from group "G"`)
}

//go:embed test/unknown_external_group.yaml
var UnknownExternalGroupConfig string

func TestParseUnknownExternalGroup(t *testing.T) {
	reader := strings.NewReader(UnknownExternalGroupConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `unknown external group "Missing" in "Go" transition from group "G"`)
}

//go:embed test/external_state.yaml
var ExternalStateConfig string

func TestParseExternalState(t *testing.T) {
	reader := strings.NewReader(ExternalStateConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, `unknown external group "A"`)
}

//go:embed test/boolean_external.yaml
var BooleanExternalConfig string

func TestParseBooleanExternal(t *testing.T) {
	reader := strings.NewReader(BooleanExternalConfig)

	_, err := parser.Parse(reader)

	require.ErrorContains(t, err, "cannot unmarshal")
}

//go:embed test/group_membership_order.yaml
var GroupMembershipOrderConfig string

func TestGroupMembershipFollowsDeclarationOrder(t *testing.T) {
	reader := strings.NewReader(GroupMembershipOrderConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, []string{"Outer", "Left", "Right", "Leaf"}, machine.StatesMap["A"].Groups.Values())
}

//go:embed test/overlapping_group_membership.yaml
var OverlappingGroupMembershipConfig string

func TestOverlappingGroupMembership(t *testing.T) {
	reader := strings.NewReader(OverlappingGroupMembershipConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, []string{"Left"}, machine.StatesMap["A"].Groups.Values())
	assert.Equal(t, []string{"Left", "Right"}, machine.StatesMap["B"].Groups.Values())
	assert.Equal(t, []string{"Right"}, machine.StatesMap["C"].Groups.Values())
}

//go:embed test/state_outside_groups.yaml
var StateOutsideGroupsConfig string

func TestStateOutsideGroupsHasNoMembership(t *testing.T) {
	reader := strings.NewReader(StateOutsideGroupsConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Empty(t, machine.StatesMap["Outside"].Groups.Values())
}

//go:embed test/initial_subgroup_chain.yaml
var InitialSubgroupChainConfig string

func TestInitialSubgroupChain(t *testing.T) {
	reader := strings.NewReader(InitialSubgroupChainConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, "A", machine.GroupsMap["Outer"].Initial.Name)
}

//go:embed test/untargeted_group_without_initial.yaml
var UntargetedGroupWithoutInitialConfig string

func TestUntargetedGroupNeedsNoInitial(t *testing.T) {
	reader := strings.NewReader(UntargetedGroupWithoutInitialConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Nil(t, machine.GroupsMap["G"].Initial)
}

//go:embed test/enter_group_initial.yaml
var EnterGroupInitialConfig string

func TestEnteringGroupSelectsInitialState(t *testing.T) {
	reader := strings.NewReader(EnterGroupInitialConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	require.Len(t, machine.StatesMap["Outside"].Transitions["Enter"], 1)
	enter := machine.StatesMap["Outside"].Transitions["Enter"][0]
	assert.Equal(t, "A", enter.Destination.Name)
	assert.Equal(t, []string{"Outer", "Left", "Right", "Leaf"}, groupNames(enter.EnterGroups))
	assert.Empty(t, enter.ExitGroups)
}

//go:embed test/enter_overlapping_group.yaml
var EnterOverlappingGroupConfig string

func TestEnteringOverlappingGroupRetainsDestination(t *testing.T) {
	reader := strings.NewReader(EnterOverlappingGroupConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	require.Len(t, machine.StatesMap["A"].Transitions["Move"], 1)
	move := machine.StatesMap["A"].Transitions["Move"][0]
	assert.Equal(t, "B", move.Destination.Name)
	assert.Equal(t, []string{"Overlap"}, groupNames(move.EnterGroups))
	assert.Empty(t, move.ExitGroups)
}

//go:embed test/state_rule_overrides_group.yaml
var StateRuleOverridesGroupConfig string

func TestUnguardedStateRuleOverridesGroupFallback(t *testing.T) {
	reader := strings.NewReader(StateRuleOverridesGroupConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	moves := machine.StatesMap["A"].Transitions["Move"]
	require.Len(t, moves, 1)
	assert.Equal(t, "B", moves[0].Destination.Name)
}

//go:embed test/inherited_group_rule_deduplication.yaml
var InheritedGroupRuleDeduplicationConfig string

func TestInheritedGroupRuleIsNotDuplicatedAcrossMembershipPaths(t *testing.T) {
	reader := strings.NewReader(InheritedGroupRuleDeduplicationConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	moves := machine.StatesMap["A"].Transitions["Move"]
	require.Len(t, moves, 1)
	assert.Equal(t, "CanMove", moves[0].Guard)
}

//go:embed test/group_transition_metadata.yaml
var GroupTransitionMetadataConfig string

func TestGroupTransitionMetadata(t *testing.T) {
	reader := strings.NewReader(GroupTransitionMetadataConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, []string{"OnMove"}, machine.ActionsMap["Move"].Transitions.Values())
	moves := machine.StatesMap["A"].Transitions["Move"]
	require.Len(t, moves, 1)
	assert.Equal(t, "OnMove", moves[0].Transition)
}

//go:embed test/action_metadata_order.yaml
var ActionMetadataOrderConfig string

func TestActionMetadataFollowsStateThenGroupDeclarationOrder(t *testing.T) {
	reader := strings.NewReader(ActionMetadataOrderConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, []string{"StateBGuard", "StateAGuard", "GroupZGuard", "GroupYGuard"}, machine.ActionsMap["Go"].Guards.Values())
	assert.Equal(t, []string{"StateBHandler", "StateAHandler", "GroupZHandler", "GroupYHandler"}, machine.ActionsMap["Go"].Transitions.Values())
}

//go:embed test/shadowed_rule_metadata.yaml
var ShadowedRuleMetadataConfig string

func TestActionMetadataIncludesShadowedRules(t *testing.T) {
	reader := strings.NewReader(ShadowedRuleMetadataConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, []string{"CanMove"}, machine.ActionsMap["Go"].Guards.Values())
	assert.Equal(t, []string{"OnMove"}, machine.ActionsMap["Go"].Transitions.Values())
}

//go:embed test/empty_group_metadata.yaml
var EmptyGroupMetadataConfig string

func TestActionMetadataIncludesGroupsWithoutMembers(t *testing.T) {
	reader := strings.NewReader(EmptyGroupMetadataConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, []string{"CanMove"}, machine.ActionsMap["Go"].Guards.Values())
	assert.Equal(t, []string{"OnMove"}, machine.ActionsMap["Go"].Transitions.Values())
}

//go:embed test/transitive_initial_group.yaml
var TransitiveInitialGroupConfig string

func TestTransitiveInitialGroup(t *testing.T) {
	reader := strings.NewReader(TransitiveInitialGroupConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, "A", machine.GroupsMap["Outer"].Initial.Name)
}

//go:embed test/transitive_initial_state.yaml
var TransitiveInitialStateConfig string

func TestTransitiveInitialState(t *testing.T) {
	reader := strings.NewReader(TransitiveInitialStateConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	assert.Equal(t, "B", machine.GroupsMap["Outer"].Initial.Name)
}

//go:embed test/external_nested_groups.yaml
var ExternalNestedGroupsConfig string

func TestExternalNestedGroupsExitAndEnterInDeclarationOrder(t *testing.T) {
	reader := strings.NewReader(ExternalNestedGroupsConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	require.Len(t, machine.StatesMap["A"].Transitions["Reset"], 1)
	reset := machine.StatesMap["A"].Transitions["Reset"][0]
	assert.Equal(t, []string{"Leaf", "Right", "Left", "Outer"}, groupNames(reset.ExitGroups))
	assert.Equal(t, []string{"Outer", "Left", "Right", "Leaf"}, groupNames(reset.EnterGroups))
}

//go:embed test/external_overlapping_group.yaml
var ExternalOverlappingGroupConfig string

func TestExternalOverlappingGroupRetainsOtherMemberships(t *testing.T) {
	reader := strings.NewReader(ExternalOverlappingGroupConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	require.Len(t, machine.StatesMap["B"].Transitions["Reset"], 1)
	reset := machine.StatesMap["B"].Transitions["Reset"][0]
	assert.Equal(t, []string{"Overlap"}, groupNames(reset.ExitGroups))
	assert.Equal(t, []string{"Overlap"}, groupNames(reset.EnterGroups))
}

//go:embed test/external_subgroup.yaml
var ExternalSubgroupConfig string

func TestExternalSubgroupDoesNotReenterParentsOrEquivalentGroups(t *testing.T) {
	reader := strings.NewReader(ExternalSubgroupConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	require.Len(t, machine.StatesMap["A"].Transitions["Reset"], 1)
	reset := machine.StatesMap["A"].Transitions["Reset"][0]
	assert.Equal(t, []string{"Leaf", "Left"}, groupNames(reset.ExitGroups))
	assert.Equal(t, []string{"Left", "Leaf"}, groupNames(reset.EnterGroups))
}

//go:embed test/empty_external_list.yaml
var EmptyExternalListConfig string

func TestEmptyExternalListRetainsSharedGroups(t *testing.T) {
	reader := strings.NewReader(EmptyExternalListConfig)

	machine, err := parser.Parse(reader)

	require.NoError(t, err)
	require.Len(t, machine.StatesMap["A"].Transitions["Move"], 1)
	move := machine.StatesMap["A"].Transitions["Move"][0]
	assert.Empty(t, move.ExitGroups)
	assert.Empty(t, move.EnterGroups)
}

func groupNames(groups []*config.Group) []string {
	names := make([]string, 0, len(groups))
	for _, group := range groups {
		names = append(names, group.Name)
	}
	return names
}
